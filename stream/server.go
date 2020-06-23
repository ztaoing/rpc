package main

import (
	"flag"
	"google.golang.org/grpc"
	"log"
	"net"
	stream_pb "rpc/stream-pb"
	"rpc/stream/service"
)

func main() {
	flag.Parse()
	listen, err := net.Listen("tcp", "127.0.0.1:1234")
	if err != nil {
		log.Fatalf("failed to listen:%v", err)
	}
	grpcServer := grpc.NewServer()
	stringService := new(service.StringService)
	stream_pb.RegisterStreamServiceServer(grpcServer, stringService)
	grpcServer.Serve(listen)
}
