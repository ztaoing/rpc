package service

import (
	"context"
	"io"
	"log"
	stream_pb "rpc/stream-pb"
	"strings"
)

const StrMaxSize = 1024

type StringService struct {
}

func (s *StringService) Concat(ctx context.Context, req *stream_pb.StringRequest) (*stream_pb.StringResponse, error) {
	if len(req.A)+len(req.B) > StrMaxSize {
		response := stream_pb.StringResponse{Res: ""}
		return &response, nil
	}
	response := stream_pb.StringResponse{Res: req.B + req.A}
	return &response, nil
}

//服务端流调用的服务端代码
func (s *StringService) LotsOfServerStream(req *stream_pb.StringRequest, qs stream_pb.StreamService_LotsOfServerStreamServer) error {
	response := stream_pb.StringResponse{Res: req.A + req.B}
	for i := 0; i < 10; i++ {
		qs.Send(&response)
	}
	return nil
}

//服务端流调用的客户端代码
func (s *StringService) LotsOfClientStream(qs stream_pb.StreamService_LotsOfClientStreamServer) error {
	var params []string
	//一直接收
	for {
		in, err := qs.Recv()
		if err == io.EOF {
			qs.SendAndClose(&stream_pb.StringResponse{Res: strings.Join(params, "")})
			return nil
		}
		if err != nil {
			log.Printf("failed to recive:%v", err)
			return err
		}
		//没有出现错误的时候
		params = append(params, in.A, in.B)

	}
}

//双向流
func (s *StringService) LotsOfServerAndClientStream(qs stream_pb.StreamService_LotsOfServerAndClientStreamServer) error {
	for {
		in, err := qs.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Printf("failed to recv:%v", err)
		}
		qs.Send(&stream_pb.StringResponse{Res: in.A + in.B})
	}
	return nil
}
