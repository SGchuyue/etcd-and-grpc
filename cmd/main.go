package main

import (
	"github.com/SGchuyue/logger/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	r "watch_etcd/server"
)

type service struct{}

func main() {
	logger.InitLogger("test.log", 1, 1, 7, false)
	// 监听本地的8080端口
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		logger.Error("监听时出现问题", err)
		return
	}
	gs := grpc.NewServer() // 创建gRPC服务器
	//	pb.Register(gs, &service{}) // 在gRPC服务端注册服务
	reflection.Register(gs) //在给定的gRPC服务器上注册服务器反射服务
	// Serve方法在lis上接受传入连接，为每个连接创建一个ServerTransport和server的goroutine。
	// 该goroutine读取gRPC请求，然后调用已注册的处理程序来响应它们。
	if err := gs.Serve(lis); err != nil {
		logger.Error("调用服务发生错误:", err)
	}
	// etcd服务注册
	etcd := r.NewService(r.Info{
		Name: "etcdtest",
		IP:   "127.0.0.1:645",
		Type: "DEBUG",
	}, []string{"127.0.0.1:4339", "127.0.0.1:2299", "127.0.0.1:1776"}) // etcd的节点ip
	etcd.Run()
	gk := etcd.GetValue()
	go etcd.Watch(gk)
}
