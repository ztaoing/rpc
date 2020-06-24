package string_service

import (
	"context"
	"github.com/go-kit/kit/log"
	"time"
)

type loggingMiddleware struct {
	Service
	logger log.Logger
}

func LoggingMidlleware(logger log.Logger) ServiceMiddleware {
	return func(service Service) Service {
		return loggingMiddleware{service, logger}
	}
}

func (mw loggingMiddleware) Concat(ctx context.Context, a, b string) (res string, err error) {
	//执行完 打印日志
	defer func(begin time.Time) {
		mw.logger.Log(
			"function", "Concat",
			"a", a,
			"b", b,
			"result", res,
			"took", time.Since(begin),
		)
	}(time.Now())
	//调用service层
	res, err = mw.Service.Concat(ctx, a, b)
	return res, err
}

func (mw loggingMiddleware) Diff(ctx context.Context, a, b string) (res string, err error) {
	//执行完 打印日志
	defer func(begin time.Time) {
		mw.logger.Log(
			"function", "Diff",
			"a", a,
			"b", b,
			"result", res,
			"took", time.Since(begin),
		)
	}(time.Now())
	//调用service层
	res, err = mw.Service.Diff(ctx, a, b)
	return res, err
}

func (mw loggingMiddleware) HealthCheck(ctx context.Context) (result bool, err error) {
	//执行完 打印日志
	defer func(begin time.Time) {
		mw.logger.Log(
			"function", "HealthCheck",
			"result", result,
			"took", time.Since(begin),
		)
	}(time.Now())
	//获得健康检查结果
	result, err = mw.Service.Health(ctx)
	return result, err
}
