package main

import (
	"context"
	"flag"
	"github.com/go-kit/kit/log"
	"google.golang.org/grpc"
	"net"
	"os"
	service "rpc/go-kit/string-service"
	"rpc/pb"
)

//在main函数中进行组装，
//首先创建StringServer，然后调用grpc_transport的NewServer方法，传入对应的endpoint和编码解码函数，得到对应的处理器，
//并赋值给StringService,然后调用gRPC的NewServer方法，并并赋值给StringService进行注册，成功启动grpc服务器

func main() {
	flag.Parse()
	ctx := context.Background()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	var svc service.Service
	svc = service.StringService{}

	//添加日志中间件
	svc = service.LoggingMidlleware(logger)(svc)

	//string endpoint
	endpoint := service.MakeStringEndpoint(svc)
	//健康检查 endpoint
	healthEndpoint := service.MakeHealthCheckEndpoint(svc)

	//封装StringEndpoint
	endpts := service.StringEndpoints{
		StringEndpoint: endpoint,
		HealthEndpoint: healthEndpoint,
	}

	handler := service.NewStringServer(ctx, endpts)
	l, _ := net.Listen("tcp", "127.0.0.1:8800")

	gRpcServer := grpc.NewServer()
	pb.RegisterStringServiceServer(gRpcServer, handler)
	gRpcServer.Serve(l)

}
