package main

import (
	"log"
	"net"

	"github.com/Pilesoft/gochat/chat"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:9000")
	if err != nil {
		log.Fatalf("Error listenting on TCP port: %v", err)
	}

	s := chat.NewServer()

	grpcServer := grpc.NewServer()

	chat.RegisterChatServiceServer(grpcServer, s)

	log.Println("Starting gRPC server on port 9000...")
	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf("Error serving gRPC: %v", err)
	}
}
