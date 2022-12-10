package main

import (
	pb "github.com/Yosorable/ms-shared/protoc_gen/admin"
	mgrpc "github.com/Yosorable/ms-shared/utils/grpc"

	"github.com/Yosorable/ms-admin/global"
	"github.com/Yosorable/ms-admin/init_service"
)

func main() {
	grpcServer := init_service.InitService()

	pb.RegisterAdminServer(grpcServer, &adminServer{})

	mgrpc.RunRpcServerInLocalHost(
		global.CONFIG.ServiceName,
		grpcServer,
	)
}
