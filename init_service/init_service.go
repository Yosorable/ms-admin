package init_service

import (
	"log"
	"runtime/debug"

	"github.com/Yosorable/ms-shared/utils"
	"github.com/Yosorable/ms-shared/utils/database"
	"google.golang.org/grpc"

	"github.com/Yosorable/ms-admin/core/handler"
	"github.com/Yosorable/ms-admin/global"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
)

func InitService() *grpc.Server {
	// 读取配置
	if err := utils.LoadJsonConfigFileWithDefaultPath(&global.CONFIG); err != nil {
		panic(err)
	}

	conf := global.CONFIG

	if db, err := database.InitMysql(
		conf.MySQL.User,
		conf.MySQL.Password,
		conf.MySQL.Addr,
		conf.MySQL.DBName,
	); err != nil {
		panic(err)
	} else {
		global.DATABASE = db
	}

	return grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(
				grpc_recovery.WithRecoveryHandler(func(p any) (err error) {
					log.Printf("[panic] %v\n %s\n", p, debug.Stack())
					return handler.NewStatusError(p)
				}),
			),
		)),
	)
}
