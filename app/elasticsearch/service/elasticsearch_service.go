package service

import (
	"coen-chat/app/elasticsearch/model"
	"coen-chat/app/elasticsearch/repository"

	"github.com/elastic/go-elasticsearch/v8"
)

type ElasticsearchService struct {
	Repository *repository.ElasticsearchRepository
}

func NewElasticsearchService(elasticsearchClient *elasticsearch.Client) *ElasticsearchService {
	return &ElasticsearchService{
		Repository: repository.NewElasticsearchRepository(elasticsearchClient),
	}
}

func (s ElasticsearchService) CreateMessage(model *model.ElasticsearchModel) {
	s.Repository.CreateMessage(model)
}

func (s ElasticsearchService) UpdateMessage(model *model.ElasticsearchModel) {
	s.Repository.UpdateMessage(model)
}

func (s ElasticsearchService) DeleteMessage(model *model.ElasticsearchModel) {
	s.Repository.DeleteMessage(model)
}
