/**
* @Author:zhoutao
* @Date:2020/7/25 下午5:57
 */

package main

import (
	"flag"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	"rpc/proto"
	"rpc/server"
)

var port string //公用端口

//同端口 不同方法实现双流支持

func main() {
	flag.StringVar(&port, "port", "8003", "多协议启动端口")
	flag.Parse()

	//初始化tcp listener ，因为grpc(http2)和http1.1在网络层上都是基于tcp协议的
	l, err := RunTcpServer(port)
	if err != nil {
		log.Fatalf("RunTcpServer err:%v", err)
	}
	//使用cmux 实现对多种协议的支持
	m := cmux.New(l)
	//application/grpc标识，grpc有特定的标识，cmux也是基于这个标识进行分流的
	grpcL := m.MatchWithWriters(cmux.HTTP2MatchHeaderFieldPrefixSendSettings("content-type", "application/grpc"))

	httpL := m.Match(cmux.HTTP1Fast())

	grpcS := RunGrpcServer(port)
	httpS := RunHttpServer(port)

	go grpcS.Serve(grpcL)
	go httpS.Serve(httpL)

	err = m.Serve()
	if err != nil {
		log.Fatalf("Run Serve err:%v", err)
	}
}

func RunTcpServer(port string) (net.Listener, error) {
	return net.Listen("tcp", ":"+port)
}

func RunHttpServer(port string) *http.Server {
	//ping路由可用于基本的心跳检测
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/ping", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("pong"))
	})
	return &http.Server{
		Addr:    ":" + port,
		Handler: serveMux,
	}
}

func RunGrpcServer(port string) *grpc.Server {
	s := grpc.NewServer()
	proto.RegisterTagServiceServer(s, server.NewTagServer())
	//使用grpcurl的前提是grpc server已经注册了反射服务
	reflection.Register(s)

	return s
}
