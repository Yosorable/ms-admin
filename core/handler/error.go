package handler

import "google.golang.org/grpc/status"

var (
	ErrorUsernameOrPassword = status.Error(7001, "用户名或密码错误")
	ErrorRegister           = status.Error(7002, "注册失败")
	ErrorToken              = status.Error(7003, "请重新登录")
)

func NewStatusError(err error) error {
	return status.Error(7000, err.Error())
}
