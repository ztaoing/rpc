/**
* @Author:zhoutao
* @Date:2020/7/26 下午7:35
 */

package middleware

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"rpc/global"
	"rpc/pkg/metatext"
	"time"
)

//超时控制
func defaultContextTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	var cancel context.CancelFunc

	//若未设置截止时间则返回false
	if _, ok := ctx.Deadline(); !ok {
		defaultTimeout := time.Second * 60
		//设置超时时间
		ctx, cancel = context.WithTimeout(ctx, defaultTimeout)
	}
	return ctx, cancel
}

//一元方法超时控制
func UnaryContextTimeout() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx, cancel := defaultContextTimeout(ctx)
		if cancel != nil {
			defer cancel()
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

//流式超时控制
func StreamContextTimeout() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		ctx, cancel := defaultContextTimeout(ctx)
		if cancel != nil {
			defer cancel()
		}
		return streamer(ctx, desc, cc, method, opts...)
	}
}

//客户端链路追踪
func ClientTracing() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		var parentContext opentracing.SpanContext
		var spanOpts []opentracing.StartSpanOption
		//解析上下文，检查是否包含上一级的span信息
		var parentSpan = opentracing.SpanFromContext(ctx)

		if parentSpan != nil {
			parentContext = parentSpan.Context()
			spanOpts = append(spanOpts, opentracing.ChildOf(parentContext))
		}

		spanOpts = append(spanOpts, []opentracing.StartSpanOption{
			opentracing.Tag{Key: string(ext.Component), Value: "grpc"},
			ext.SpanKindRPCClient,
		}...)

		span := global.Tracer.StartSpan(method, spanOpts...)
		defer span.Finish()

		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}
		_ = global.Tracer.Inject(span.Context(), opentracing.TextMap, metatext.MetadataTextMap{md})
		newContext := opentracing.ContextWithSpan(metadata.NewOutgoingContext(ctx, md), span)
		return invoker(newContext, method, req, reply, cc, opts...)
	}
}
