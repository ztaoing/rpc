package main

import (
	"flag"
	"google.golang.org/grpc"
	"log"
	"net"
	string_service "rpc/grpc/string-service"
	"rpc/pb"
)

func main() {
	flag.Parse()
	listen, err := net.Listen("tcp", "127.0.0.1:1234")
	if err != nil {
		log.Fatalf("failed to listen:%v", err)
	}
	grpcServer := grpc.NewServer()
	stringService := new(string_service.StringService)
	pb.RegisterStringServiceServer(grpcServer, stringService)
	grpcServer.Serve(listen)
}
