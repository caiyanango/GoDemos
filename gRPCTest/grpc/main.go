package main

import (
	"context"
	"fmt"
	"gRPCTest/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

type greeterServer struct {
	pb.UnimplementedGreeterServer
}

func (s *greeterServer) SayHello(c context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Reply: "message: " + req.Name}, nil
}

func main() {
	listener, _ := net.Listen("tcp", ":8000")
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &greeterServer{})
	reflection.Register(s)
	fmt.Println("服务中")
	s.Serve(listener)
}
