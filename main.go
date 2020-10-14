package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:9000")
	if err != nil {
		log.Fatalf("Error listenting on TCP port: %v", err)
	}

	grpcServer := grpc.NewServer()

	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf("Error serving gRPC: %v", err)
	}
}
