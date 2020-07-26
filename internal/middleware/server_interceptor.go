/**
* @Author:zhoutao
* @Date:2020/7/26 下午5:18
 */

package middleware

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"rpc/pkg/errcode"
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
