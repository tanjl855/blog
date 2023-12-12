package main

import (
	"context"
	"fmt"

	"github.com/tanjl855/blog/internal/blog/proto"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("start to dial grpc server")
	innerIp := "localhost:8084"
	conn, err := grpc.Dial(innerIp, grpc.WithDefaultCallOptions(
		grpc.MaxCallRecvMsgSize(1024*1024*1000), // 最大消息接收大小为100MB
		grpc.MaxCallSendMsgSize(1024*1024*1000), // 最大消息发送大小为100MB
	), grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
		return
	}
	cli := proto.NewUserServiceClient(conn)
	rsp, err := cli.UserTest(context.Background(), &proto.User{
		Username: "tanjl test",
		Age:      23,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(rsp)
}
