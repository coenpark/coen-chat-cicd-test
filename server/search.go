package server

import (
	"context"
	"log"
	"net"
	"os"

	"coen-chat/app/search/service"
	searchpb "coen-chat/protos/search"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func SearchGRPCServerHandler(mux *runtime.ServeMux, esClient *elasticsearch.Client) {
	searchConn := SearchGRPCServer(esClient)
	err := searchpb.RegisterSearchServiceHandler(context.Background(), mux, searchConn)
	if err != nil {
		log.Fatalln("Failed register user gateway", err)
	}
}

func SearchGRPCServer(esClient *elasticsearch.Client) *grpc.ClientConn {
	socket := os.Getenv("GRPC_SEARCH_SOCKET")
	lis, err := net.Listen("tcp", socket)
	if err != nil {
		panic(err)
	}
	log.Println("Listening search grpc at", socket)

	s := grpc.NewServer()
	s.RegisterService(&searchpb.SearchService_ServiceDesc, service.NewSearchService(esClient))
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
