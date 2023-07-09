package service

import (
	"context"
	"errors"
	"regexp"

	"coen-chat/app/user/model"
	"coen-chat/app/user/repository"
	userpb "coen-chat/protos/user"
	"coen-chat/utils"

	"github.com/go-redis/redis/v7"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type UserService struct {
	Repository *repository.UserRepository
	userpb.UserServiceServer
}

func NewUserService(db *gorm.DB, redis *redis.Client) userpb.UserServiceServer {
	return &UserService{
		Repository: repository.NewUserRepository(db, redis),
	}
}

func (s *UserService) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error) {
	// validate email
	if match, _ := regexp.MatchString("^[\\w-\\.]+@([\\w-]+\\.)+[\\w-]{2,4}$", req.GetEmail()); !match {
		return nil, errors.New("잘못된 이메일 유형입니다")
	}
	if req.GetName() == "" || req.GetEmail() == "" || req.GetPassword() == "" {
		return nil, errors.New("데이터를 모두 입력해주세요")
	}
	// 기존 회원가입 여부 확인
	if s.Repository.IsExist(req.GetEmail()) {
		return nil, errors.New("이미 존재하는 아이디입니다")
	}
	encodedPassword, err := HashPassword(req.GetPassword())
	if err != nil {
		return nil, errors.New("비밀번호 해싱 실패")
	}
	// 회원 객체 생성
	user := &model.User{
		Email:    req.GetEmail(),
		Password: encodedPassword,
		Name:     req.GetName(),
	}
	// 회원 등록
	s.Repository.CreateUser(user)
	// JWT 생성
	token, err := utils.CreateJWTToken(user)
	if err != nil {
		return nil, errors.New("토큰 생성 실패")
	}
	err = s.Repository.CreateAuth(user.Email, token)
	if err != nil {
		return nil, errors.New("토큰 redis 저장 실패")
	}
	return &userpb.CreateUserResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}, nil
}

func (s *UserService) Login(ctx context.Context, req *userpb.LoginRequest) (*userpb.LoginResponse, error) {
	if req.GetEmail() == "" || req.GetPassword() == "" {
		return nil, errors.New("데이터를 모두 입력해주세요")
	}
	user, err := s.Repository.FindByEmail(req.GetEmail())
	if err != nil {
		return nil, errors.New("일치하는 사용자가 없습니다")
	}
	if !CheckPasswordHash(req.GetPassword(), user.Password) {
		return nil, errors.New("비밀번호가 일치하지 않습니다")
	}
	token, err := utils.CreateJWTToken(&user)
	if err != nil {
		return nil, errors.New("토큰 생성 실패")
	}
	err = s.Repository.CreateAuth(user.Email, token)
	if err != nil {
		return nil, errors.New("토큰 redis 저장 실패")
	}
	return &userpb.LoginResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}, nil
}

func (s *UserService) Logout(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	accessToken, refreshToken := utils.ExtractTokenFromContext(ctx)
	at, _ := utils.ExtractTokenMetadata(accessToken, true)
	s.Repository.DeleteAuth(at.AccessUUID)
	rt, _ := utils.ExtractTokenMetadata(refreshToken, false)
	s.Repository.DeleteAuth(rt.AccessUUID)

	return &emptypb.Empty{}, nil
}

func HashPassword(password string) (string, error) {
	cost := 14
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
