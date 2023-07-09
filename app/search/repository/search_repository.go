package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"

	searchpb "coen-chat/protos/search"

	"github.com/elastic/go-elasticsearch/v8"
)

type SearchRepository struct {
	esClient *elasticsearch.Client
}

func NewSearchRepository(esClient *elasticsearch.Client) *SearchRepository {
	return &SearchRepository{
		esClient: esClient,
	}
}

func (s *SearchRepository) SearchMessage(from, to, keyword string) ([]*searchpb.Message, error) {
	var model []*searchpb.Message
	var err error
	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"match": map[string]interface{}{
							"content": keyword,
						},
					},
					{
						"range": map[string]interface{}{
							"createdAt": map[string]interface{}{
								"gte": from,
								"lte": to,
							},
						},
					},
				},
			},
		},
	}

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}
	res, err := s.esClient.Search(
		s.esClient.Search.WithContext(context.Background()),
		s.esClient.Search.WithIndex("message"),
		s.esClient.Search.WithBody(&buf),
		s.esClient.Search.WithTrackTotalHits(true),
		s.esClient.Search.WithPretty(),
	)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Fatalf("Error parsing the response body: %s", err)
		} else {
			// Print the response status and error information.
			log.Fatalf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	// Print the response status, number of results, and request duration.
	log.Printf(
		"[%s] %d hits; took: %dms",
		res.Status(),
		int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)),
		int(r["took"].(float64)),
	)
	// Print the ID and document source for each hit.
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		var mod searchpb.Message
		res := hit.(map[string]interface{})["_source"].(map[string]interface{})
		mod.MessageId = res["messageId"].(string)
		mod.Email = res["email"].(string)
		mod.Content = res["content"].(string)
		mod.ChannelId = fmt.Sprintf("%d", int(res["channelId"].(float64)))
		mod.CreatedAt = res["createdAt"].(string)
		mod.UpdatedAt = res["updatedAt"].(string)
		model = append(model, &mod)
	}
	return model, err
}
