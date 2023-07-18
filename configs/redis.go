package configs

import (
	"log"
	"os"

	"github.com/go-redis/redis/v7"
)

func ConnectRedis() *redis.Client {
	dsn := os.Getenv("GCP_REDIS_DSN")
	password := os.Getenv("REDIS_PASSWORD")
	client := redis.NewClient(&redis.Options{
		Addr:     dsn,
		Password: password,
		DB:       0,
	})
	_, err := client.Ping().Result()
	if err != nil {
		log.Fatalln("Failed to connect Redis", err)
	}
	log.Println("redis connected")
	return client
}
