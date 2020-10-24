package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	pb "watch_etcd/proto"
	s "watch_etcd/server"
)

type service struct{}

func (s *service) SendTest(ctx context.Context, req *pb.TestRequest) (res *pb.TestResponse, err error) {
	fmt.Printf("邮箱:%s;发送内容:%s", req.Send, req.Text)
	return &pb.TestResponse{
		Ok: true,
	}, nil
}

func main() {
	// 监听本地的8972端口
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
		return
	}
	gs := grpc.NewServer() // 创建gRPC服务器
	//	pb.Register(gs, &service{}) // 在gRPC服务端注册服务
	reflection.Register(gs) //在给定的gRPC服务器上注册服务器反射服务
	// Serve方法在lis上接受传入连接，为每个连接创建一个ServerTransport和server的goroutine。
	// 该goroutine读取gRPC请求，然后调用已注册的处理程序来响应它们。
	//etcd服务注册
	reg := s.NewService(s.Info{
		Name: "etcdtest",
		IP:   "127.0.0.1:645", // grpc服务节点ip
	}, []string{"127.0.0.1:2379", "127.0.0.1:2279", "127.0.0.1:3379"}) // etcd的节点ip
	go reg.Run()
	gk := reg.GetValue()
	go reg.Watch(gk)
	if err := gs.Serve(lis); err != nil {
		fmt.Println(err)
	}
	fmt.Println(222)
}
