/**
* @Author:zhoutao
* @Date:2020/7/26 上午7:08
 */

package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"rpc/proto"
)

func main() {
	ctx := context.Background()
	clientConn, err := GetClientConn(ctx, "localhost:8004", nil)
	if err != nil {
		log.Fatalf("err:", err)
	}
	defer clientConn.Close()

	//初始化client对象
	tagServiceClient := proto.NewTagServiceClient(clientConn)
	//发起rpc方法调用
	resp, err := tagServiceClient.GetTagList(ctx, &proto.GetTagListRequest{
		Name: "go",
	})
	if err != nil {
		log.Fatalf("tagServiceClient.GetTagList err:%v", err)
	}
	log.Printf("resp: %v", resp)

}

func GetClientConn(ctx context.Context, target string, opts []grpc.DialOption) (*grpc.ClientConn, error) {
	opts = append(opts, grpc.WithInsecure())
	//在opts中加入withBlock模式，这样发起拨号连接时会阻塞等待连接完成，是最终连接到达ready状态，这是连接才正式可用
	//grpc.DialContext是异步建立连接的，并不会马上就成为可用连接，仅处于connection状态（需要多久取决于外部因素，如网络）
	//只有正式到达ready状态这个链接才算真正的可用
	return grpc.DialContext(ctx, target, opts...)
}
