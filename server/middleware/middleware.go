package middleware

import (
	"context"
	"errors"
	"log"
	"net/http"

	"coen-chat/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"google.golang.org/protobuf/proto"
)

type UnsignedResponse struct {
	Message interface{} `json:"message"`
}

type TokenExpired struct {
	Message interface{} `json:"message"`
	Code    interface{} `json:"code"`
}

func ForwardResponseOption(c context.Context, w http.ResponseWriter, msg proto.Message) error {
	headers := w.Header()
	if location, ok := headers["Grpc-Metadata-Location"]; ok {
		w.Header().Set("Location", location[0])
		w.WriteHeader(http.StatusFound)
	}
	return nil
}

func ValidateTokenMiddleware(c *gin.Context) {
	// Request Header 에서 데이터 추출
	var accessToken string
	var refreshToken string
	var redisClient *redis.Client
	// gin context에서 redis client 꺼내오기
	redisCtx, _ := c.Get("redisClient")
	redisClient = redisCtx.(*redis.Client)
	authorization, err := utils.ExtractTokenMetadata(c.GetHeader("authorization"), true)
	if err != nil {
		log.Println("metadata 추출 실패", err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, UnsignedResponse{Message: errors.New("토큰인증 실패").Error()})
		return
	}
	refreshUUID, err := utils.ExtractTokenMetadata(c.GetHeader("refreshToken"), false)
	if err != nil {
		log.Println("metadata 추출 실패", err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, UnsignedResponse{Message: errors.New("토큰인증 실패").Error()})
		return
	}
	accessToken = c.GetHeader("authorization")
	refreshToken = c.GetHeader("refreshToken")
	// Redis AccessUUID 확인
	_, err = utils.FetchAuth(authorization, redisClient)
	if err != nil { // access token uuid 없으면
		_, err := utils.FetchAuth(refreshUUID, redisClient) // refresh token 확인
		if err != nil {
			log.Println("middleware : refresh token uuid redis에 없음", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, TokenExpired{Code: 401, Message: errors.New("토큰인증 실패").Error()})
			return
		}
		// 새로 만들어서 보내줌
		refresh, err := utils.Refresh(c.GetHeader("refreshToken"))
		if err != nil {
			log.Println("refresh token 으로 새로운 토큰 생성 실패", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, TokenExpired{Code: 401, Message: errors.New("토큰인증 실패").Error()})
			return
		}
		accessToken = refresh.AccessToken
		refreshToken = refresh.RefreshToken
	}
	c.Request.Header.Set("Grpc-Metadata-access-token", accessToken)
	c.Request.Header.Set("Grpc-Metadata-refresh-token", refreshToken)
}
