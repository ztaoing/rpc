/**
* @Author:zhoutao
* @Date:2020/7/25 下午4:50
 */

package server

import (
	"context"
	"encoding/json"
	"rpc/pkg/bapi"
	"rpc/pkg/errcode"
	"rpc/proto"
)

type TagServer struct {
}

func NewTagServer() *TagServer {
	return &TagServer{}
}

func (t *TagServer) GetTagList(ctx context.Context, r *proto.GetTagListRequest) (*proto.GetTagListReply, error) {
	//博客后端服务地址
	api := bapi.NewAPI("http://127.0.0.1:8000")
	body, err := api.GetTagList(ctx, r.GetName())
	if err != nil {
		return nil, errcode.TogRPCError(errcode.ERROR_GET_TAG_LIST_FAIL)
	}

	tagList := proto.GetTagListReply{}
	err = json.Unmarshal(body, &tagList)
	if err != nil {
		return nil, err
	}
	return &tagList, nil
}
