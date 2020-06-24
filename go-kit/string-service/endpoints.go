package string_service

import (
	"context"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"rpc/pb"
	"strings"
)

//endpoint层主要负责request/response格式的转换，以及公用拦截器相关的逻辑。
//是go-kit的核心，采用洋葱模式，提供了对日志、限流、熔断、链路追踪和服务监控等方面的扩展能力

type StringEndpoints struct {
	StringEndpoint endpoint.Endpoint
	HealthEndpoint endpoint.Endpoint
}

func (ue *StringEndpoints) Health(ctx context.Context) (bool, error) {
	panic("implement me")
}

/*func (ue *StringEndpoints) Health(ctx context.Context) (bool, error) {
	status, err := ue.HealthEndpoint(ctx)
	response := status.(*HealthResponse)
	return response.Status, err
}*/

func (ue *StringEndpoints) Concat(ctx context.Context, a, b string) (string, error) {
	resp, err := ue.StringEndpoint(ctx, &pb.StringRequest{A: a, B: b})
	response := resp.(*pb.StringResponse)
	return response.Ret, err
}

func (ue *StringEndpoints) Diff(ctx context.Context, a, b string) (string, error) {
	resp, err := ue.StringEndpoint(ctx, &pb.StringRequest{A: a, B: b})
	response := resp.(*pb.StringResponse)
	return response.Ret, err
}

//string request/response
type StringRequest struct {
	RequestType string `json:"request_type"`
	A           string `json:"a"`
	B           string `json:"b"`
}

type StringResponse struct {
	Result string `json:"result"`
	Error  error  `json:"error"`
}

var ErrInvalidRequestType = errors.New("invalid request type,options:Concat or Diff")

func MakeStringEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(StringRequest)
		var (
			res, a, b string
			opError   error
		)
		a = req.A
		b = req.B
		if strings.EqualFold(req.RequestType, "Concat") {
			res, opError = svc.Concat(ctx, a, b)
		} else if strings.EqualFold(req.RequestType, "Diff") {
			res, opError = svc.Diff(ctx, a, b)
		} else {
			return nil, ErrInvalidRequestType
		}
		return StringResponse{Result: res, Error: opError}, nil

	}
}

//health request/response
type HealthRequest struct {
}

type HealthResponse struct {
	Status bool `json:"status"`
}

func MakeHealthCheckEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		//获得状态...
		status := true
		return HealthResponse{
			Status: status,
		}, nil

	}
}
