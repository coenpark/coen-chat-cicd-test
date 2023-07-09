package server

import (
	"context"
	"log"
	"net"
	"os"

	"coen-chat/app/message/service"
	messagepb "coen-chat/protos/message"

	"cloud.google.com/go/pubsub"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/gorm"
)

func MessageGRPCServerHandler(mux *runtime.ServeMux, db *gorm.DB, pubsubClient *pubsub.Client) {
	messageConn := MessageGRPCServer(db, pubsubClient)
	err := messagepb.RegisterMessageServiceHandler(context.Background(), mux, messageConn)
	if err != nil {
		log.Fatalln("Failed register user gateway", err)
	}
}

func MessageGRPCServer(db *gorm.DB, pubsubClient *pubsub.Client) *grpc.ClientConn {
	socket := os.Getenv("GRPC_MESSAGE_SOCKET")
	lis, err := net.Listen("tcp", socket)
	if err != nil {
		panic(err)
	}
	log.Println("Listening message grpc at", socket)

	s := grpc.NewServer()
	s.RegisterService(&messagepb.MessageService_ServiceDesc, service.NewMessageService(db, pubsubClient))
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
