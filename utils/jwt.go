package utils

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"coen-chat/app/user/model"
	"coen-chat/configs"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v7"
	"github.com/google/uuid"
)

type AuthTokenClaims struct {
	Authorized bool
	AccessUUID string
	Email      string
	Exp        int64

	jwt.StandardClaims
}

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUUID   string
	RefreshUUID  string
	AtExpires    int64
	RtExpires    int64
}

type AccessDetails struct {
	AccessUUID string
	Email      string
}

var redisClient *redis.Client

func init() {
	redisClient = configs.ConnectRedis()
}

func CreateJWTToken(user *model.User) (*TokenDetails, error) {
	td := &TokenDetails{}
	// AccessToken
	random, err := uuid.NewRandom()
	expTime := 5
	td.AtExpires = time.Now().Add(time.Duration(expTime) * time.Minute).Unix()
	td.AccessUUID = random.String()
	if err != nil {
		log.Println("uuid create fail")
		return nil, err
	}
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["email"] = user.Email
	atClaims["access_uuid"] = td.AccessUUID
	atClaims["name"] = user.Name
	atClaims["exp"] = time.Now().Add(time.Duration(expTime) * time.Minute)
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("JWT_ACCESS_TOKEN_SECRET")))
	if err != nil {
		return nil, err
	}

	// RefreshToken
	refreshExpTime := 1
	td.RtExpires = time.Now().Add(time.Hour * time.Duration(expTime)).Unix()
	random, err = uuid.NewRandom()
	if err != nil {
		log.Println("uuid create fail")
		return nil, err
	}
	td.RefreshUUID = random.String()
	rtClaim := jwt.MapClaims{}
	rtClaim["refresh_uuid"] = td.RefreshUUID
	rtClaim["email"] = user.Email
	rtClaim["name"] = user.Name
	rtClaim["exp"] = time.Now().Add(time.Duration(refreshExpTime) * time.Hour)
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaim)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("JWT_REFRESH_TOKEN_SECRET")))
	return td, nil
}

func VerifyJWTToken(tokenString string, isAccessToken bool) (*jwt.Token, error) {
	var tokenSecret string = ""
	if isAccessToken {
		tokenSecret = os.Getenv("JWT_ACCESS_TOKEN_SECRET")
	} else {
		tokenSecret = os.Getenv("JWT_REFRESH_TOKEN_SECRET")
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func ExtractTokenMetadata(token string, isAccessToken bool) (*AccessDetails, error) {
	var claimUUID = ""
	if isAccessToken {
		claimUUID = "access_uuid"
	} else {
		claimUUID = "refresh_uuid"
	}
	tok, err := VerifyJWTToken(token, isAccessToken)
	if err != nil {
		return nil, err
	}
	claims, ok := tok.Claims.(jwt.MapClaims)
	if ok && tok.Valid {
		accessUUID, ok := claims[claimUUID].(string)
		if !ok {
			return nil, err
		}
		email, ok := claims["email"].(string)
		if !ok {
			return nil, err
		}
		return &AccessDetails{
			AccessUUID: accessUUID,
			Email:      email,
		}, nil
	}
	return nil, err
}

func FetchAuth(authD *AccessDetails, redisClient *redis.Client) (string, error) {
	email, err := redisClient.Get(authD.AccessUUID).Result()
	if err != nil {
		return "", err
	}
	return email, nil
}

func Refresh(refreshToken string) (*TokenDetails, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_REFRESH_TOKEN_SECRET")), nil
	})
	if err != nil {
		return nil, errors.New("parse Refresh Token Fail")
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return nil, errors.New("refresh Token expired")
	}
	claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
	if ok && token.Valid {
		refreshUuid, ok := claims["refresh_uuid"].(string) //convert the interface to string
		if !ok {
			return nil, errors.New("refreshUuid not found")
		}
		email, ok := claims["email"].(string)
		if !ok {
			return nil, errors.New("email not found")
		}
		name, ok := claims["name"].(string)
		if !ok {
			return nil, errors.New("name not found")
		}
		//Delete the previous Refresh Token

		deleted, err := redisClient.Del(refreshUuid).Result()
		if err != nil || deleted == 0 {
			return nil, errors.New("redis delete fail")
		}
		//Create new pairs of refresh and access tokens
		td, createErr := CreateJWTToken(&model.User{Email: email, Name: name})
		if createErr != nil {
			return nil, errors.New("토큰 생성 실패")
		}
		//save the tokens metadata to redis
		at := time.Unix(td.AtExpires, 0) //converting Unix to UTC
		rt := time.Unix(td.RtExpires, 0)
		now := time.Now()

		errAccess := redisClient.Set(td.AccessUUID, email, at.Sub(now)).Err()
		if errAccess != nil {
			return nil, errAccess
		}
		errRefresh := redisClient.Set(td.RefreshUUID, email, rt.Sub(now)).Err()
		if errRefresh != nil {
			return nil, errRefresh
		}
		return td, nil
	} else {
		return nil, errors.New("refresh token expired")
	}
}
