package init_service

import (
	"github.com/Yosorable/ms-shared/utils"
	"github.com/Yosorable/ms-shared/utils/database"

	"github.com/Yosorable/ms-admin/global"
)

func InitService() {
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
}
