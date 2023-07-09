package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"coen-chat/app/elasticsearch/model"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type ElasticsearchRepository struct {
	elasticsearchClient *elasticsearch.Client
}

func NewElasticsearchRepository(elasticsearchClient *elasticsearch.Client) *ElasticsearchRepository {
	return &ElasticsearchRepository{
		elasticsearchClient: elasticsearchClient,
	}
}

func (r *ElasticsearchRepository) CreateMessage(message *model.ElasticsearchModel) {
	fmt.Println(message)
	// Build the request body.
	data, err := json.Marshal(message)
	if err != nil {
		log.Fatalf("Error marshaling document: %s", err)
	}

	// Set up the request object.
	req := esapi.IndexRequest{
		Index:      "message",
		DocumentID: message.MessageId,
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	// Perform the request with the client.
	res, err := req.Do(context.Background(), r.elasticsearchClient)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

}

func (r *ElasticsearchRepository) UpdateMessage(message *model.ElasticsearchModel) {
	// Build the request body.
	data, err := json.Marshal(message)
	if err != nil {
		log.Fatalf("Error marshaling document: %s", err)
	}

	// Set up the request object.
	req := esapi.IndexRequest{
		Index:      "message",
		DocumentID: message.MessageId,
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	// Perform the request with the client.
	res, err := req.Do(context.Background(), r.elasticsearchClient)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()
}

func (r *ElasticsearchRepository) DeleteMessage(message *model.ElasticsearchModel) {
	fmt.Println(message)
	// Set up the request object.
	req := esapi.DeleteRequest{
		Index:      "message",
		DocumentID: message.MessageId,
		Refresh:    "true",
	}

	// Perform the request with the client.
	res, err := req.Do(context.Background(), r.elasticsearchClient)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()
}
