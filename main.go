package main

import (
	"context"
	"log"
	"net/http"
	"os"

	esService "coen-chat/app/elasticsearch/service"
	"coen-chat/app/pubsub/service"
	"coen-chat/configs"
	"coen-chat/server"
	"coen-chat/server/middleware"
	"coen-chat/utils"

	"cloud.google.com/go/pubsub"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"gorm.io/gorm"
)

var db *gorm.DB
var redisClient *redis.Client
var elasticsearchDB *elasticsearch.Client
var pubsubClient *pubsub.Client

func init() {
	configs.LoadEnv()
	db = configs.ConnectMysql()
	elasticsearchDB = configs.ConnectElasticsearch()
	pubsubClient = configs.NewPubsubClient()
	redisClient = configs.ConnectRedis()
}

func main() {
	// DB 테이블 초기화
	//db.AutoMigrate(&userModel.User{}, &channelModel.Channel{}, &channelModel.Member{}, &messageModel.Message{})

	// Socket.io 초기화
	serveMux := http.NewServeMux()
	go server.StartSocketServer(serveMux)

	gwmux := runtime.NewServeMux(
		runtime.WithForwardResponseOption(middleware.ForwardResponseOption),
	)

	esService := esService.NewElasticsearchService(elasticsearchDB)
	server.UserGRPCServerHandler(gwmux, db, redisClient)
	server.ChannelGRPCServerHandler(gwmux, db)
	server.MessageGRPCServerHandler(gwmux, db, pubsubClient)
	server.SearchGRPCServerHandler(gwmux, elasticsearchDB)
	go service.Subscribe(context.Background(), pubsubClient, esService)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("redisClient", redisClient)
		c.Next()
	})
	cors := utils.GetCors()
	r.Use(cors)
	r.Use(gin.Logger())

	path := r.Group("/api/v1")
	path.Use(middleware.ValidateTokenMiddleware)
	{
		path.Any("/*path", gin.WrapF(gwmux.ServeHTTP))
	}
	socketPath := r.Group("/socket.io")
	socketPath.Use(cors)
	{
		socketPath.GET("/*any", gin.WrapH(serveMux))
		socketPath.POST("/*any", gin.WrapH(serveMux))
	}
	authPath := r.Group("/api")
	{
		authPath.POST("/join", gin.WrapF(gwmux.ServeHTTP))
		authPath.POST("/login", gin.WrapF(gwmux.ServeHTTP))
	}
	//r.Run(os.Getenv("GIN_GONIC_SOCKET"))
	log.Println("CICD TEST LOG")
	r.Run(os.Getenv("GIN_GONIC_SOCKET"))
	log.Println("Serving gRPC-Gateway on ", os.Getenv("GATEWAY_SOCKET"))
}
