package main

import (
	"context"
	"flag"
	"fmt"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
	service "rpc/go-kit/string-service"
	"rpc/pb"
)

func main() {
	flag.Parse()
	ctx := context.Background()
	conn, err := grpc.Dial("127.0.0.1:8800", grpc.WithInsecure())
	if err != nil {
		fmt.Println("grpc dail err:", err)
	}
	defer conn.Close()

	server := NewStringClient(conn)
	result, err := server.Concat(ctx, "A", "B")
	if err != nil {
		fmt.Println("check error", err.Error())
	}
	fmt.Println("result = ", result)
}

func NewStringClient(conn *grpc.ClientConn) service.Service {
	//构建client
	var ep = grpctransport.NewClient(conn,
		"pb.StringService",
		"Concat",
		EncodeStringRequest,
		DecodeStringResponse,
		pb.StringResponse{},
	).Endpoint()
	userEp := service.StringEndpoints{
		StringEndpoint: ep,
	}
	return &userEp
}

func EncodeStringRequest(ctx context.Context, r interface{}) (interface{}, error) {
	return r, nil
}

func DecodeStringResponse(ctx context.Context, r interface{}) (interface{}, error) {
	return r, nil
}
