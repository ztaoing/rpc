/**
* @Author:zhoutao
* @Date:2020/7/25 下午5:57
 */

package main

import (
	"flag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	"rpc/proto"
	"rpc/server"
)

var grpcPort string
var httpPort string

func main() {
	flag.StringVar(&grpcPort, "grpc_port", "8001", "grpc启动端口号")
	flag.StringVar(&httpPort, "http_port", "9001", "HTTP启动端口号")
	flag.Parse()

	errs := make(chan error)
	//为什么放在goroutine中？是因为HTTPEndPoint和GRPCEndPoint实际上是阻塞行为
	go func() {
		err := RunHttpServer(httpPort)
		if err != nil {
			errs <- err
		}
	}()

	go func() {
		err := RunGrpcServer(grpcPort)
		if err != nil {
			errs <- err
		}
	}()

	select {
	case err := <-errs:
		log.Fatalf("Run Server err:", err)
	}
}

func RunHttpServer(port string) error {
	//ping路由可用于基本的心跳检测
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/ping", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("pong"))
	})
	return http.ListenAndServe(":"+port, serveMux)
}

func RunGrpcServer(port string) error {
	s := grpc.NewServer()
	proto.RegisterTagServiceServer(s, server.NewTagServer())
	//使用grpcurl的前提是grpc server已经注册了反射服务
	reflection.Register(s)
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}
	return s.Serve(lis)
}
