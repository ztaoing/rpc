package string_service

import (
	"context"
	"github.com/go-kit/kit/transport/grpc"
	"rpc/pb"
)

/**
transport层主要负责网络传输，处理HTTP、grpc、thrift等相关逻辑
*/
//定义grpcServer结构体，他有两个grpc.handler方法，分别是concat和diff

type grpcServer struct {
	concat grpc.Handler
	diff   grpc.Handler
}

func (s *grpcServer) Concat(ctx context.Context, r *pb.StringRequest) (*pb.StringResponse, error) {
	_, resp, err := s.concat.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.StringResponse), nil
}

func (s *grpcServer) Diff(ctx context.Context, r *pb.StringRequest) (*pb.StringResponse, error) {
	_, resp, err := s.diff.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.StringResponse), nil
}

func NewStringServer(ctx context.Context, endpoint StringEndpoints) pb.StringServiceServer {
	return &grpcServer{
		concat: grpc.NewServer(
			endpoint.StringEndpoint,
			DecodeConcatStringRequest,
			EncodeStringResponse,
		),
		diff: grpc.NewServer(
			endpoint.StringEndpoint,
			DecodeDiffStringRequest,
			EncodeStringResponse,
		),
	}
}

//定义encode和decode
func DecodeConcatStringRequest(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.StringRequest)
	return StringRequest{
		RequestType: "Concat",
		A:           req.A,
		B:           req.B,
	}, nil
}

func DecodeDiffStringRequest(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.StringRequest)
	return StringRequest{
		RequestType: "Diff",
		A:           req.A,
		B:           req.B,
	}, nil
}

func EncodeStringResponse(ctx context.Context, r interface{}) (interface{}, error) {
	resp := r.(StringResponse)
	if resp.Error != nil {
		return &pb.StringResponse{
			Ret: resp.Result,
			Err: resp.Error.Error(),
		}, nil
	}
	//未出错时
	return &pb.StringResponse{
		Ret: resp.Result,
		Err: "",
	}, nil
}
