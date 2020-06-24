package string_service

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

//service层是执行业务逻辑的地方
const StrMaxSize = 1024

var (
	ErrMaxSize = errors.New("maximum size of 1024 bytes exceeded")
	ErrStrSize = errors.New("maximum size of 1024 bytes exceeded")
)

type Service interface {
	Concat(ctx context.Context, a, b string) (string, error)
	Diff(ctx context.Context, a, b string) (string, error)
	Health(ctx context.Context) (bool, error)
}

type StringService struct {
}

func (s StringService) Concat(ctx context.Context, a, b string) (string, error) {
	if len(a)+len(b) > StrMaxSize {
		return "", ErrMaxSize
	}
	fmt.Printf("StringService Concat: %s concat %s = %s\n", a, b, a+b)
	return a + b, nil
}

func (s StringService) Diff(ctx context.Context, a, b string) (string, error) {
	if len(a) < 1 || len(b) < 1 {
		return "", nil
	}
	res := ""
	if len(a) > len(b) {
		for _, char := range b {
			if strings.Contains(a, string(char)) {
				res = res + string(char)
			}
		}
	} else {
		for _, char := range a {
			if strings.Contains(b, string(char)) {
				res = res + string(char)
			}
		}
	}
	return res, nil
}

func (s StringService) Health(ctx context.Context) (bool, error) {
	status := true
	fmt.Printf("health check")
	return status, nil
}

//定义中间件
type ServiceMiddleware func(Service) Service
