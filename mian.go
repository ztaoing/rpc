/**
* @Author:zhoutao
* @Date:2020/7/25 下午5:57
 */

package main

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"rpc/proto"
	"rpc/server"
)

var port = "8001"

func main() {
	s := grpc.NewServer()
	proto.RegisterTagServiceServer(s, server.NewTagServer())
	//使用grpcurl的前提是grpc server已经注册了反射服务
	reflection.Register(s)

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("net.Listen err:", err)
	}

	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("server.serve err:", err)
	}
}
