package main

import (
	"context"
	"fmt"
	"net"

	"github.com/tanjl855/blog/internal/blog/proto"
	"google.golang.org/grpc"
)

func main() {
	// 模拟其他服务器上面实现的不同逻辑或者其他方法，通过内网ip去调用
	go func() {
		rpcServer := grpc.NewServer()

		proto.RegisterUserServiceServer(rpcServer, &UserServiceInSecondServer{})

		listener, err := net.Listen("tcp", ":8084")
		if err != nil {
			panic(err)
		}

		fmt.Println("rpc server start at 8084")
		err = rpcServer.Serve(listener)
		if err != nil {
			panic(err)
		}
	}()
	rpcServer := grpc.NewServer()

	proto.RegisterUserServiceServer(rpcServer, &UserService{})

	listener, err := net.Listen("tcp", ":8082")
	if err != nil {
		panic(err)
	}

	fmt.Println("rpc server start at 8082")
	err = rpcServer.Serve(listener)
	if err != nil {
		panic(err)
	}
}

var _ = proto.UserServiceServer(&UserService{})

type UserService struct {
}

func (u *UserService) UserTest(ctx context.Context, user *proto.User) (*proto.ErrRsp, error) {
	// logic
	fmt.Println(user)
	return &proto.ErrRsp{
		ErrCode: 0,
		ErrMsg:  "test",
	}, nil
}

type UserServiceInSecondServer struct {
}

func (u *UserServiceInSecondServer) UserTest(ctx context.Context, user *proto.User) (*proto.ErrRsp, error) {
	fmt.Println("other server")
	fmt.Println(user)
	return &proto.ErrRsp{
		ErrCode: 0,
		ErrMsg:  "test in other server",
	}, nil
}
