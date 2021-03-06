/**
* @Author:zhoutao
* @Date:2020/7/26 下午1:20
 */

package main

import (
	"context"
	"encoding/json"
	"flag"
	assetfs "github.com/elazarl/go-bindata-assetfs"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"log"
	"net/http"
	"path"
	"rpc/internal/middleware"
	"rpc/pkg/swagger"
	"rpc/proto"
	"rpc/server"
	"strings"
)

var port string

func init() {
	flag.StringVar(&port, "port", "8004", "启动端口")
	flag.Parse()
}

func main() {
	err := RunServer(port)
	if err != nil {
		log.Fatalf("Run server err:%v", err)
	}
}

//不同协议的分流
func grpcHandlerFunc(grpcServer *grpc.Server, otherHanler http.Handler) http.Handler {
	//h2c.NewHandler其内部逻辑是拦截所有的h2c流量，根据不同的请求流量类型将其劫持并重定向到响应的handler中去处理
	//最终完成在同个端口上既能提供http/1.1的功能，又能提供http/2的功能
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//ProtoMajor是客户端请求的版本号，客户端始终是使用http/1.1或者http/2
		//Content-Type 确定流量的类型
		//h2c标识允许通过明文tcp运行http/2协议，此标识符用于http/1.1升级标头字段和标识http/2 over tcp
		//而官方标准库golang.org/x/net/http2/h2c实现了http/2的未加密模式，直接使用即可
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHanler.ServeHTTP(w, r)
		}
	}), &http2.Server{})
}

func RunServer(port string) error {
	httpMux := runHttpServer()

	gatewayMux := runGrpcGatewayServer()

	grpcS := runGrpcServer()

	//注册grpc gateway的错误方法
	runtime.HTTPError = grpcGatewayError

	httpMux.Handle("/", gatewayMux)

	return http.ListenAndServe(":"+port, grpcHandlerFunc(grpcS, httpMux))
}

func runHttpServer() *http.ServeMux {
	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/ping", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("pong"))
	})

	//swagger
	prefix := "/swagger-ui/"
	fileServer := http.FileServer(&assetfs.AssetFS{
		//由于引用的原因，实际上这两处是没有问题的
		Asset:    swagger.Asset,
		AssetDir: swagger.AssetDir,
		Prefix:   "third_party/swagger-ui",
	})
	serverMux.Handle(prefix, http.StripPrefix(prefix, fileServer))
	serverMux.HandleFunc("/swagger/", func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "swagger.json") {
			http.NotFound(w, r)
			return
		}
		p := strings.TrimPrefix(r.URL.Path, "/swagger/")
		p = path.Join("proto", p)
		http.ServeFile(w, r, p)
	})
	return serverMux
}

func runGrpcServer() *grpc.Server {
	//grpc server的相关属性都可以在ServerOption中设置，如keep-alive、credentials等参数
	opts := []grpc.ServerOption{
		//grpc_middleware是通过递归interceptor数组的方式执行代表rpc方法的handler.
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			HelloInterceptor,
			WorldInterceptor,
			middleware.ServerTracing, //链路追踪
		)),
	}
	s := grpc.NewServer(opts...)
	//添加拦截器

	proto.RegisterTagServiceServer(s, server.NewTagServer())
	reflection.Register(s)
	return s
}

//server端拦截
func HelloInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Println("你好")
	resp, err := handler(ctx, req)
	log.Println("再见")
	return resp, err
}

//server端拦截器
func WorldInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Println("你好，红烧煎鱼")
	resp, err := handler(ctx, req)
	log.Println("再见")
	return resp, err
}

func runGrpcGatewayServer() *runtime.ServeMux {
	endpoint := "0.0.0.0:" + port
	gatewayMux := runtime.NewServeMux()
	//指定为非加密模式
	//grpc server/client在启动和调用时，必须明确其是否加密
	opts := []grpc.DialOption{grpc.WithInsecure()}
	//注册tagServiceHandler事件,其内部会自动转换并拨号到grpc endpoint，并在上下文结束后关闭连接
	proto.RegisterTagServiceHandlerFromEndpoint(context.Background(), gatewayMux, endpoint, opts)
	return gatewayMux
}

type httpError struct {
	Code    int32  `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

//对返回的grpc错误进行了两次处理，将其转化为对应的HTTP状态码和对应的消息主体，以确保客户端能够根据restful api的标准进行交互
func grpcGatewayError(ctx context.Context, _ *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, _ *http.Request, err error) {
	s, ok := status.FromError(err)
	if !ok {
		s = status.New(codes.Unknown, err.Error())
	}
	httpError := httpError{Code: int32(s.Code()), Message: s.Message()}

	details := s.Details()
	for _, detail := range details {
		if v, ok := detail.(*proto.Error); ok {
			httpError.Code = v.Code
			httpError.Message = v.Message
		}
	}
	resp, _ := json.Marshal(httpError)
	w.Header().Set("Content-type", marshaler.ContentType())
	w.WriteHeader(runtime.HTTPStatusFromCode(s.Code()))
	w.Write(resp)
}
