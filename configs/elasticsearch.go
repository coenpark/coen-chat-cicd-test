package configs

import (
	"log"
	"os"

	"github.com/elastic/go-elasticsearch/v8"
)

func ConnectElasticsearch() *elasticsearch.Client {
	log.SetFlags(0)

	config := elasticsearch.Config{
		CloudID: os.Getenv("ELS_CLOUD_ID"),
		APIKey:  os.Getenv("ELS_API_KEY"),
	}
	es, err := elasticsearch.NewClient(config)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	return es
}
