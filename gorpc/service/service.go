package service

import (
	"errors"
	"strings"
)

const StrMaxSize = 1024

var (
	ErrMaxSize = errors.New("maximum size of 1024 bytes exceeded")
)

//定义远程过程调用相关接口传入参数和返回参数的数据结构
type StringRequest struct {
	A string
	B string
}

//定义一个服务对象
type Service interface {
	Concat(req StringRequest, ret *string) error
	Diff(req StringRequest, ret *string) error
}

//实现了Service接口
type StringService struct {
}

func (s StringService) Concat(req StringRequest, ret *string) error {
	//是否溢出
	if len(req.A)+len(req.B) > StrMaxSize {
		*ret = ""
		return ErrMaxSize
	}
	*ret = req.A + req.B
	return nil
}

func (s StringService) Diff(req StringRequest, ret *string) error {
	if len(req.A) < 1 || len(req.B) < 1 {
		*ret = ""
		return nil
	}
	res := ""
	if len(req.A) >= len(req.B) {
		for _, char := range req.B {
			if strings.Contains(req.B, string(char)) {
				res = res + string(char)
			}
		}
	} else {
		for _, char := range req.A {
			if strings.Contains(req.A, string(char)) {
				res = res + string(char)
			}
		}
	}
	*ret = res
	return nil
}

type ServiceMiddleware func(Service) Service
