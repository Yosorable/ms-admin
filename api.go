package main

import (
	"context"

	pb "github.com/Yosorable/ms-shared/protoc_gen/admin"

	"github.com/Yosorable/ms-admin/core/handler"
)

type server struct {
	pb.UnimplementedAdminServer
}

func (*server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginReply, error) {
	return handler.Login(ctx, req)
}

func (*server) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterReply, error) {
	return handler.Register(ctx, req)
}

func (*server) CheckToken(ctx context.Context, req *pb.CheckTokenRequest) (*pb.CheckTokenReply, error) {
	return handler.CheckToken(ctx, req)
}

func (*server) GetUserByID(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.GetUserByIDReply, error) {
	return handler.GetUserByID(ctx, req)
}
