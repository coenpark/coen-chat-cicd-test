package server

import (
	"context"
	"log"
	"net"
	"os"

	"coen-chat/app/channel/service"
	channelpb "coen-chat/protos/channel"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/gorm"
)

func ChannelGRPCServerHandler(mux *runtime.ServeMux, db *gorm.DB) {
	channelConn := ChannelGRPCServer(db)
	err := channelpb.RegisterChannelServiceHandler(context.Background(), mux, channelConn)
	if err != nil {
		log.Fatalln("Failed register channel gateway", err)
	}
}

func ChannelGRPCServer(db *gorm.DB) *grpc.ClientConn {
	socket := os.Getenv("GRPC_CHANNEL_SOCKET")
	lis, err := net.Listen("tcp", socket)
	if err != nil {
		panic(err)
	}
	log.Println("Listening channel grpc at", socket)

	s := grpc.NewServer()
	s.RegisterService(&channelpb.ChannelService_ServiceDesc, service.NewChannelService(db))
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
