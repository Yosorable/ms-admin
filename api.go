package main

import (
	"context"

	pb "github.com/Yosorable/ms-shared/protoc_gen/admin"

	"github.com/Yosorable/ms-admin/core/handler"
)

type adminServer struct {
	pb.UnimplementedAdminServer
}

func (*adminServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginReply, error) {
	return handler.Login(ctx, req)
}

func (*adminServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterReply, error) {
	return handler.Register(ctx, req)
}

func (*adminServer) CheckToken(ctx context.Context, req *pb.CheckTokenRequest) (*pb.CheckTokenReply, error) {
	return handler.CheckToken(ctx, req)
}

func (*adminServer) GetUserByID(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.GetUserByIDReply, error) {
	return handler.GetUserByID(ctx, req)
}

func (*adminServer) CreateUserRecordTableIfNotExist(ctx context.Context, req *pb.CreateUserRecordTableIfNotExistRequest) (*pb.CreateUserRecordTableIfNotExistReply, error) {
	return handler.CreateUserRecordTableIfNotExist(ctx, req)
}

func (*adminServer) QueryUserRecord(ctx context.Context, req *pb.QueryUserRecordRequest) (*pb.QueryUserRecordReply, error) {
	return handler.QueryUserRecord(ctx, req)
}

func (*adminServer) CreateOrUpdateUserRecord(ctx context.Context, req *pb.CreateOrUpdateUserRecordRequest) (*pb.CreateOrUpdateUserRecordReply, error) {
	return handler.CreateOrUpdateUserRecord(ctx, req)
}

func (*adminServer) DeleteUserRecord(ctx context.Context, req *pb.DeleteUserRecordRequest) (*pb.DeleteUserRecordReply, error) {
	return handler.DeleteUserRecord(ctx, req)
}
