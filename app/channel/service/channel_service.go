package service

import (
	"context"
	"errors"
	"log"

	"coen-chat/app/channel/model"
	"coen-chat/app/channel/repository"
	channelpb "coen-chat/protos/channel"
	"coen-chat/utils"

	"gorm.io/gorm"
)

type ChannelService struct {
	Repository *repository.ChannelRepository
	channelpb.ChannelServiceServer
}

func NewChannelService(db *gorm.DB) channelpb.ChannelServiceServer {
	return &ChannelService{
		Repository: repository.NewChannelRepository(db),
	}
}

func (s *ChannelService) GetChannel(ctx context.Context, req *channelpb.GetChannelRequest) (*channelpb.GetChannelResponse, error) {
	accessToken, refreshToken := utils.ExtractTokenFromContext(ctx)
	channel, err := s.Repository.FindByIdPreload(req.GetChannelId())
	if err != nil {
		log.Println("채널 조회 실패")
		return nil, errors.New("채널 조회 실패")
	}

	var members []*channelpb.Member
	for _, member := range channel.Members {
		toProto := &channelpb.Member{
			ChannelId: member.ChannelID,
			Email:     member.Email,
		}
		members = append(members, toProto)
	}
	var messages []*channelpb.Message
	for _, message := range channel.Messages {
		toProto := &channelpb.Message{
			MessageId: message.MessageID,
			Email:     message.Email,
			Content:   message.Content,
			UpdatedAt: message.UpdatedAt.String(),
			CreatedAt: message.CreatedAt.String(),
		}
		messages = append(messages, toProto)
	}

	responseModel := &channelpb.GetChannelResponse{
		Channel: &channelpb.Channel{
			ChannelId:   channel.ChannelID,
			ChannelName: channel.ChannelName,
			CreatedBy:   channel.CreatedBy,
			UpdatedAt:   channel.UpdatedAt.String(),
		},
		TokenResponse: &channelpb.TokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
		Member:  members,
		Message: messages,
	}
	return responseModel, nil
}

func (s *ChannelService) GetChannelList(ctx context.Context, req *channelpb.GetChannelListRequest) (*channelpb.GetChannelListResponse, error) {
	accessToken, refreshToken := utils.ExtractTokenFromContext(ctx)
	tokenMetadata, err := utils.ExtractTokenMetadata(accessToken, true)
	if err != nil {
		log.Println("토큰 메타데이터 추출 실패", err)
		return nil, errors.New("토큰 메타데이터 추출 실패")
	}
	channelList, err := s.Repository.GetChannelList(req.GetIsJoined(), tokenMetadata.Email)
	if err != nil {
		log.Println("채널 리스트 조회 실패", err)
		return nil, errors.New("채널 리스트 조회 실패")
	}
	var response []*channelpb.Channel
	for _, channel := range channelList {
		toProto := channelpb.Channel{
			ChannelId:   channel.ChannelID,
			ChannelName: channel.ChannelName,
			CreatedBy:   channel.CreatedBy,
			UpdatedAt:   channel.UpdatedAt.String(),
		}
		response = append(response, &toProto)
	}
	return &channelpb.GetChannelListResponse{
		Channel: response,
		TokenResponse: &channelpb.TokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}, nil
}

func (s *ChannelService) CreateChannel(ctx context.Context, req *channelpb.CreateChannelRequest) (*channelpb.CreateChannelResponse, error) {
	accessToken, refreshToken := utils.ExtractTokenFromContext(ctx)
	tokenMetadata, err := utils.ExtractTokenMetadata(accessToken, true)
	if err != nil {
		log.Println("토큰 메타데이터 추출 실패", err)
		return nil, errors.New("토큰 메타데이터 추출 실패")
	}
	var channel = &model.Channel{
		ChannelName: req.GetChannelName(),
		CreatedBy:   tokenMetadata.Email,
	}

	err = s.Repository.CreateChannel(channel)
	if err != nil {
		log.Println(err)
		return nil, errors.New("채널 생성 실패")
	}

	var member = &model.Member{
		ChannelID: channel.ChannelID,
		Email:     channel.CreatedBy,
	}
	err = s.Repository.JoinChannel(member)
	if err != nil {
		log.Println("채널 참가 실패", err)
		return nil, errors.New("채널 참가 실패")
	}
	return &channelpb.CreateChannelResponse{
		TokenResponse: &channelpb.TokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}, nil
}

func (s *ChannelService) UpdateChannel(ctx context.Context, req *channelpb.UpdateChannelRequest) (*channelpb.UpdateChannelResponse, error) {
	channel, err := s.Repository.FindById(req.GetChannelId())
	if err != nil {
		log.Println("채널 조회 실패.", err)
		return nil, errors.New("채널 조회 실패")
	}
	// jwt email, created_by 일치 확인
	accessToken, refreshToken := utils.ExtractTokenFromContext(ctx)
	tokenMetadata, err := utils.ExtractTokenMetadata(accessToken, true)
	if err != nil {
		log.Println("토큰 메타데이터 추출 실패", err)
		return nil, errors.New("토큰 메타데이터 추출 실패")
	}
	if tokenMetadata.Email != channel.CreatedBy {
		log.Println("다른 사용자 접근 :", errors.New("권한 없는 채널 수정 요청"))
		return nil, errors.New("권한이 없습니다")
	}
	channel.ChannelName = req.GetChannelName()
	s.Repository.UpdateChannel(channel)
	return &channelpb.UpdateChannelResponse{
		TokenResponse: &channelpb.TokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}, nil
}

func (s *ChannelService) DeleteChannel(ctx context.Context, req *channelpb.DeleteChannelRequest) (*channelpb.DeleteChannelResponse, error) {
	channel, err := s.Repository.FindById(req.GetChannelId())
	if err != nil {
		log.Println("채널 조회 실패.", err)
		return nil, errors.New("채널 조회 실패")
	}
	// jwt email, created_by 일치 확인
	accessToken, refreshToken := utils.ExtractTokenFromContext(ctx)
	tokenMetadata, err := utils.ExtractTokenMetadata(accessToken, true)
	if err != nil {
		log.Println("토큰 메타데이터 추출 실패", err)
		return nil, errors.New("토큰 메타데이터 추출 실패")
	}
	if tokenMetadata.Email != channel.CreatedBy {
		log.Println("다른 사용자 접근", errors.New("권한 없는 채널 삭제 요청"))
		return nil, errors.New("권한이 없습니다")
	}

	err = s.Repository.DeleteChannel(channel)
	if err != nil {
		log.Println("채널 삭제 실패", err)
		return nil, errors.New("채널 삭제 실패")
	}
	return &channelpb.DeleteChannelResponse{
		TokenResponse: &channelpb.TokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}, nil
}

func (s *ChannelService) JoinChannel(ctx context.Context, req *channelpb.JoinChannelRequest) (*channelpb.JoinChannelResponse, error) {
	accessToken, refreshToken := utils.ExtractTokenFromContext(ctx)
	tokenMetadata, err := utils.ExtractTokenMetadata(accessToken, true)
	if err != nil {
		log.Println("토큰 메타데이터 추출 실패", err)
		return nil, errors.New("토큰 메타데이터 추출 실패")
	}
	var member = &model.Member{
		ChannelID: req.GetChannelId(),
		Email:     tokenMetadata.Email,
	}
	err = s.Repository.JoinChannel(member)
	if err != nil {
		log.Println("채널 참가 실패")
		return nil, errors.New("채널 참가 실패")
	}
	return &channelpb.JoinChannelResponse{
		TokenResponse: &channelpb.TokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}, nil
}
