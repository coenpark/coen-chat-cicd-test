package server

import (
	"context"
	"log"
	"net"
	"os"

	"coen-chat/app/user/service"
	userpb "coen-chat/protos/user"

	"github.com/go-redis/redis/v7"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/gorm"
)

func UserGRPCServerHandler(mux *runtime.ServeMux, db *gorm.DB, redisDB *redis.Client) {
	userConn := UserGRPCServer(db, redisDB)
	err := userpb.RegisterUserServiceHandler(context.Background(), mux, userConn)
	if err != nil {
		log.Fatalln("Failed register user gateway", err)
	}
}

func UserGRPCServer(db *gorm.DB, redis *redis.Client) *grpc.ClientConn {
	socket := os.Getenv("GRPC_USER_SOCKET")
	lis, err := net.Listen("tcp", socket)
	if err != nil {
		panic(err)
	}
	log.Println("Listening user grpc at", socket)

	s := grpc.NewServer()
	s.RegisterService(&userpb.UserService_ServiceDesc, service.NewUserService(db, redis))
	go s.Serve(lis)
	if err != nil {
		log.Fatalln("Errored while Serving : ", socket, err)
	}

	conn, err := grpc.DialContext(
		context.Background(),
		socket,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	return conn
}
