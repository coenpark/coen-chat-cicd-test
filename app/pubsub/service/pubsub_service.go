package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	esModel "coen-chat/app/elasticsearch/model"
	esService "coen-chat/app/elasticsearch/service"
	msgModel "coen-chat/app/message/model"

	"cloud.google.com/go/pubsub"
)

func Publish(ctx context.Context, pubsubClient *pubsub.Client, crud string, message msgModel.Message) error {
	msg, err := json.Marshal(message)
	t := pubsubClient.Topic("messages")
	result := t.Publish(ctx, &pubsub.Message{
		Data: msg,
		Attributes: map[string]string{
			"crud": crud,
		},
	})
	id, err := result.Get(ctx)
	if err != nil {
		return fmt.Errorf("pubsub: result.Get: %w", err)
	}
	log.Printf("Published a message; msg ID: %v\n", id)
	return nil
}

func Subscribe(ctx context.Context, pubsubClient *pubsub.Client, esService *esService.ElasticsearchService) {
	sub := pubsubClient.Subscription("messages-sub")

	sub.ReceiveSettings.MaxOutstandingMessages = 1000
	sub.ReceiveSettings.MaxOutstandingBytes = 1e10

	err := sub.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
		var message esModel.ElasticsearchModel
		msg := string(m.Data)
		json.Unmarshal([]byte(msg), &message)
		ElasticsearchHandler(&message, m.Attributes["crud"], esService)
		m.Ack()
	})
	if err != nil {
		log.Println(err)
	}
}

func ElasticsearchHandler(message *esModel.ElasticsearchModel, crud string, esService *esService.ElasticsearchService) {
	switch crud {
	case "create":
		esService.CreateMessage(message)
	case "update":
		esService.UpdateMessage(message)
	case "delete":
		esService.DeleteMessage(message)
	}
}
