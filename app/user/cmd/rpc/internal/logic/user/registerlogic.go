package userlogic

import (
	"context"
	"go-zero-init/app/user/models"
	"go-zero-init/common/tools"
	"go-zero-init/common/xerr"

	"go-zero-init/app/user/cmd/rpc/internal/svc"
	"go-zero-init/app/user/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *pb.RegisterReq) (*pb.RegisterResp, error) {
	// 先通过 username 查询用户是否存在
	userModel := models.NewDefaultUserModel(l.svcCtx.DBEngin)
	user, err := userModel.SearchUserByUsername(in.Username)
	// 如果存在，返回用户已经存在，注册失败
	if user != nil {
		return nil, xerr.NewErrCode(xerr.UserExistError)
	}
	if err != nil && err.(*xerr.CodeError).GetErrCode() != xerr.RecordNotFoundError {
		return nil, err
	}
	// 用户第一次注册，调用 createUser 创建用户
	// 处于数据安全考虑，用户密码存入数据库前先做加密处理
	pwd, pwdErr := encodeUserPassword(in.Password)
	if pwdErr != nil {
		return nil, pwdErr
	}
	userMap := map[string]interface{}{
		"username": in.Username,
		"password": pwd,
	}
	err = userModel.CreateUser(userMap)
	if err != nil {
		return nil, xerr.NewErrCode(xerr.CreateUserError)
	}

	return &pb.RegisterResp{Username: in.Username}, nil
}

// 用户密码加密

func encodeUserPassword(pwd string) (string, error) {
	hashStr, err := tools.EncodeBcrypt(pwd)
	if err != nil {
		return "", xerr.NewErrCode(xerr.EncryptionError)
	}
	return tools.EncodeMd5([]byte(hashStr)), nil
}
