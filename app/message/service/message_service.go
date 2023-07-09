package service

import (
	"context"
	"errors"
	"log"
	"time"

	"coen-chat/app/message/model"
	"coen-chat/app/message/repository"
	"coen-chat/app/pubsub/service"
	messagepb "coen-chat/protos/message"
	"coen-chat/utils"

	"cloud.google.com/go/pubsub"
	"gorm.io/gorm"
)

type MessageService struct {
	Repository   *repository.MessageRepository
	PubsubClient *pubsub.Client
	messagepb.MessageServiceServer
}

func NewMessageService(db *gorm.DB, pubsubClient *pubsub.Client) messagepb.MessageServiceServer {
	return &MessageService{
		Repository:   repository.NewMessageRepository(db),
		PubsubClient: pubsubClient,
	}
}

func (s *MessageService) CreateMessage(ctx context.Context, req *messagepb.CreateMessageRequest) (*messagepb.CreateMessageResponse, error) {
	accessToken, refreshToken := utils.ExtractTokenFromContext(ctx)
	tokenMetadata, err := utils.ExtractTokenMetadata(accessToken, true)
	if err != nil {
		log.Println("토큰 메타데이터 추출 실패", err)
		return nil, errors.New("토큰 메타데이터 추출 실패")
	}
	createdAt, err := time.Parse(time.RFC3339, req.GetCreatedAt())
	if err != nil {
		log.Println("time parse fail", err)
		return nil, errors.New("시간 형식이 잘못되었습니다")
	}
	var message = &model.Message{
		MessageID: req.GetMessageId(),
		ChannelID: req.GetChannelId(),
		Email:     tokenMetadata.Email,
		Content:   req.GetContent(),
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
	}
	err = s.Repository.CreateMessage(message)
	if err != nil {
		log.Println("메세지 생성 실패", err)
		return nil, errors.New("메세지 생성 실패")
	}
	err = service.Publish(ctx, s.PubsubClient, "create", *message)
	if err != nil {
		log.Println("message_service > CreateMessage > pubsub Publish failed:", err)
	}
	return &messagepb.CreateMessageResponse{
		TokenResponse: &messagepb.TokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}, nil
}

func (s *MessageService) UpdateMessage(ctx context.Context, req *messagepb.UpdateMessageRequest) (*messagepb.UpdateMessageResponse, error) {
	accessToken, refreshToken := utils.ExtractTokenFromContext(ctx)
	tokenMetadata, err := utils.ExtractTokenMetadata(accessToken, true)
	if err != nil {
		log.Println("토큰 메타데이터 추출 실패", err)
		return nil, errors.New("토큰 메타데이터 추출 실패")
	}
	message, err := s.Repository.FindById(req.GetMessageId())
	if err != nil {
		log.Println("메세지 조회 실패", err)
		return nil, errors.New("메세지 조회 실패")
	}
	if message.Email != tokenMetadata.Email {
		log.Println("다른 사용자 접근", errors.New("권한 없는 메시지 수정 요청"))
		return nil, errors.New("권한이 없습니다")
	}
	message.Content = req.GetContent()
	err = s.Repository.UpdateMessage(message)
	if err != nil {
		log.Println("메세지 수정 실패", err)
		return nil, errors.New("메세지 수정 실패")
	}
	err = service.Publish(ctx, s.PubsubClient, "update", *message)
	if err != nil {
		log.Println("message_service > UpdateMessage > pubsub Publish failed:", err)
	}
	return &messagepb.UpdateMessageResponse{
		TokenResponse: &messagepb.TokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}, nil
}

func (s *MessageService) DeleteMessage(ctx context.Context, req *messagepb.DeleteMessageRequest) (*messagepb.DeleteMessageResponse, error) {
	accessToken, refreshToken := utils.ExtractTokenFromContext(ctx)
	tokenMetadata, err := utils.ExtractTokenMetadata(accessToken, true)
	if err != nil {
		log.Println("토큰 메타데이터 추출 실패", err)
		return nil, errors.New("토큰 메타데이터 추출 실패")
	}
	message, err := s.Repository.FindById(req.GetMessageId())
	if err != nil {
		log.Println("메세지 조회 실패", err)
		return nil, errors.New("메세지 조회 실패")
	}
	if message.Email != tokenMetadata.Email {
		log.Println("다른 사용자 접근", errors.New("권한 없는 메시지 삭제 요청"))
		return nil, errors.New("권한이 없습니다")
	}
	err = s.Repository.DeleteMessage(message)
	if err != nil {
		log.Println("메세지 삭제 실패", err)
		return nil, errors.New("메세지 삭제 실패")
	}
	err = service.Publish(ctx, s.PubsubClient, "delete", *message)
	if err != nil {
		log.Println("message_service > DeleteMessage > pubsub Publish failed:", err)
	}
	return &messagepb.DeleteMessageResponse{
		TokenResponse: &messagepb.TokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}, nil
}

func (s *MessageService) GetMessageList(ctx context.Context, req *messagepb.GetMessageListRequest) (*messagepb.GetMessageListResponse, error) {
	accessToken, refreshToken := utils.ExtractTokenFromContext(ctx)
	_, err := utils.ExtractTokenMetadata(accessToken, true)
	if err != nil {
		log.Println("토큰 메타데이터 추출 실패", err)
		return nil, errors.New("토큰 메타데이터 추출 실패")
	}
	var messages []*messagepb.Message
	messageList, err := s.Repository.GetMessageList(req.GetChannelId(), int(req.GetOffset()))

	for _, message := range messageList {
		toProto := &messagepb.Message{
			MessageId: message.MessageID,
			Email:     message.Email,
			Content:   message.Content,
			CreatedAt: message.CreatedAt.String(),
			UpdatedAt: message.UpdatedAt.String(),
		}
		messages = append(messages, toProto)
	}
	return &messagepb.GetMessageListResponse{
		Message: messages,
		TokenResponse: &messagepb.TokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}, nil
}
