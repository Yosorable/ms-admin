package handler

import (
	"fmt"

	"google.golang.org/grpc/status"
)

var (
	ErrorUsernameOrPassword = status.Error(7001, "用户名或密码错误")
	ErrorRegister           = status.Error(7002, "注册失败")
	ErrorToken              = status.Error(7003, "请重新登录")
)

func NewStatusError(err any) error {
	return status.Error(7000, fmt.Sprintf("%v", err))
}

func NotNULLError(filed string) error {
	return status.Error(7000, fmt.Sprintf("%s not allow null", filed))
}
