package configs

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
)

func NewPubsubClient() *pubsub.Client {
	ctx := context.Background()
	creds := option.WithCredentialsFile(os.Getenv("GCP_CREDENTIAL_FILE_VOLUME_PATH")) // cloud
	//creds := option.WithCredentialsFile(os.Getenv("GCP_CREDENTIAL_FILE_LOCAL_PATH")) // local
	client, err := pubsub.NewClient(ctx, os.Getenv("PUBSUB_SUB"), creds)
	if err != nil {
		log.Fatalln(err)
	}
	return client
}
