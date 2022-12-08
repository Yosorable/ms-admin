package main

import (
	pb "github.com/Yosorable/ms-shared/protoc_gen/admin"
	mgrpc "github.com/Yosorable/ms-shared/utils/grpc"

	"github.com/Yosorable/ms-admin/global"
	"github.com/Yosorable/ms-admin/init_service"

	"google.golang.org/grpc"
)

func main() {
	init_service.InitService()
	mgrpc.RunRpcServerInLocalHost(
		global.CONFIG.ServiceName,
		func(s *grpc.Server) {
			pb.RegisterAdminServer(s, &server{})
		},
	)
}
