package handler

import (
	"context"
	"time"

	pb "github.com/Yosorable/ms-shared/protoc_gen/admin"
	"github.com/Yosorable/ms-shared/protoc_gen/common"
	"github.com/Yosorable/ms-shared/utils"

	"github.com/Yosorable/ms-admin/core/model"
	admin_utils "github.com/Yosorable/ms-admin/core/utils"
	"github.com/Yosorable/ms-admin/global"
)

func Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginReply, error) {
	if req.GetUsername() == "" || req.GetPassword() == "" {
		return nil, ErrorUsernameOrPassword
	}

	var user model.User
	var cnt int64
	err := global.DATABASE.First(&user, "username = ?", req.Username).Count(&cnt).Error
	if err != nil || cnt != 1 || user.Hash != utils.Hash(req.Password+user.Salt) {
		return nil, ErrorUsernameOrPassword
	}

	jwt := admin_utils.NewJWT()
	token, err := jwt.CreateToken(jwt.CreateClaims(admin_utils.BaseClaims{
		ID:       user.ID,
		Username: user.Username,
	}))
	if err != nil {
		return nil, NewStatusError(err)
	}

	return &pb.LoginReply{
		User: &pb.User{
			Id:       int32(user.ID),
			Username: user.Username,
			TimeInfo: &common.TimeInfo{
				CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
				UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
			},
		},
		JwtToken: token,
	}, nil
}

func Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterReply, error) {
	if req.GetUsername() == "" || req.GetPassword() == "" {
		return nil, ErrorUsernameOrPassword
	}

	salt := utils.GetUUID()
	user := model.User{
		Username:  req.GetUsername(),
		Hash:      utils.Hash(req.GetPassword() + salt),
		Salt:      salt,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := global.DATABASE.Create(&user).Error
	if err != nil {
		return nil, ErrorRegister
	}
	return &pb.RegisterReply{}, nil
}

func CheckToken(ctx context.Context, req *pb.CheckTokenRequest) (*pb.CheckTokenReply, error) {
	j := admin_utils.NewJWT()
	// parseToken 解析token包含的信息
	claims, err := j.ParseToken(req.GetJwtToken())
	if err != nil {
		return nil, ErrorToken
	}
	res := &pb.CheckTokenReply{
		User: &pb.User{
			Id:       int32(claims.ID),
			Username: claims.Username,
		},
	}
	if claims.ExpiresAt-time.Now().Unix() < claims.BufferTime {
		dr, _ := admin_utils.ParseDuration(global.CONFIG.JWT.ExpiresTime)
		claims.ExpiresAt = time.Now().Add(dr).Unix()
		newToken, _ := j.CreateTokenByOldToken(req.GetJwtToken(), *claims)
		res.NewToken = newToken
	}

	return res, nil
}

func GetUserByID(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.GetUserByIDReply, error) {
	id := req.GetId()
	var user model.User

	err := global.DATABASE.First(&user, "id = ?", id).Error
	if err != nil {
		return nil, NewStatusError(err)
	}

	return &pb.GetUserByIDReply{
		User: &pb.User{
			Id:       int32(user.ID),
			Username: user.Username,
			TimeInfo: &common.TimeInfo{
				CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
				UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
			},
		},
	}, nil
}
