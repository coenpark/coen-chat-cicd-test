package service

import (
	"context"
	"errors"
	"log"

	"coen-chat/app/search/repository"
	searchpb "coen-chat/protos/search"
	"coen-chat/utils"

	"github.com/elastic/go-elasticsearch/v8"
)

type SearchService struct {
	Repository *repository.SearchRepository
	searchpb.SearchServiceServer
}

func NewSearchService(esClient *elasticsearch.Client) *SearchService {
	return &SearchService{
		Repository: repository.NewSearchRepository(esClient),
	}
}

func (s *SearchService) SearchMessage(ctx context.Context, req *searchpb.SearchMessageRequest) (*searchpb.SearchMessageResponse, error) {
	accessToken, refreshToken := utils.ExtractTokenFromContext(ctx)
	_, err := utils.ExtractTokenMetadata(accessToken, true)
	if err != nil {
		log.Println("토큰 메타데이터 추출 실패", err)
		return nil, errors.New("토큰 메타데이터 추출 실패")
	}
	from := req.GetFrom()
	to := req.GetTo()
	keyword := req.GetKeyword()
	messages, err := s.Repository.SearchMessage(from, to, keyword)
	return &searchpb.SearchMessageResponse{
		Message: messages,
		TokenResponse: &searchpb.TokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}, nil
}
