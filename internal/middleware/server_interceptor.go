/**
* @Author:zhoutao
* @Date:2020/7/26 下午5:18
 */

package middleware

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"rpc/pkg/errcode"
	"rpc/pkg/metatext"
	"runtime/debug"
	"time"
)

//访问日志
func AccessLog(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	requestLog := "access request log :method:%s,begin_time:%d,request:%v"
	begin_time := time.Now().Local().Unix()
	log.Printf(requestLog, info.FullMethod, begin_time, req)

	resp, err := handler(ctx, req)

	responseLog := "access response log:method:%s,begin_time:%d,end_time:%d,response:%v"
	end_time := time.Now().Local().Unix()
	log.Printf(responseLog, info.FullMethod, begin_time, end_time, resp)

	return resp, err
}

//错误日志
func ErrorLog(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	resp, err := handler(ctx, req)
	if err != nil {
		errLog := "error log:method:%s,code:%v,message:%v,details:%v"
		s := errcode.FromError(err)
		log.Printf(errLog, info.FullMethod, s.Code(), s.Err().Error(), s.Details())
	}
	return resp, err
}

//异常捕获
func Recovery(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	defer func() {
		if e := recover(); e != nil {
			recoverLog := "recovery log:method:%s,method:%s,message:%s,stack:%s"
			log.Printf(recoverLog, info.FullMethod, e, string(debug.Stack()[:]))
		}
	}()
	return handler(ctx, req)
}

//链路追踪
func ServerTracing(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	//从上下文中读取
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		//新生成一个空的metadata
		md = metadata.New(nil)
	}
	//从给定的载体中解码出span context实例
	parentSpanContext, _ := opentracing.Tracer.Extract(opentracing.TextMap, metatext.MetadataTextMap{md})
	//options 设置本地span的标签信息
	opts := []opentracing.StartSpanOption{
		opentracing.Tag{Key: string(ext.Component), Value: "gRPC"},
		ext.SpanKindRPCServer,
		ext.RPCServerOption(parentSpanContext),
	}

	//根据父span生成新的span
	span := opentracing.Tracer.StartSpan(info.FullMethod, opts...)
	defer span.Finish()

	ctx = opentracing.ContextWithSpan(ctx, span)
	return handler(ctx, req)
}
