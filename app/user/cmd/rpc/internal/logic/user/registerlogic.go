package userlogic

import (
	"context"
	"go-zero-init/app/user/cmd/rpc/internal/svc"
	"go-zero-init/app/user/cmd/rpc/pb"
	"go-zero-init/app/user/models/dao"
	"go-zero-init/app/user/models/entity"
	"go-zero-init/common/tools"
	"go-zero-init/common/xerr"

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
	// 1, 先通过 username 查询用户是否存在
	userModel := dao.NewDefaultUserModel(l.svcCtx.DBEngin)
	user, err := userModel.SearchUserByUsername(in.Username)
	// 如果存在，返回用户已经存在，注册失败
	if user != nil {
		return nil, xerr.NewErrCode(xerr.UserExistError)
	}
	// 出现 非记录不存在 的错误，返回错误
	if err != nil && err.(*xerr.CodeError).GetErrCode() != xerr.RecordNotFoundError {
		return nil, err
	}

	// 2, 用户第一次注册，调用 createUser 创建用户
	// 处于数据安全考虑，用户密码存入数据库前先做加密处理
	pwd, pwdErr := encodeUserPassword(in.Password)
	if pwdErr != nil {
		return nil, pwdErr
	}
	u := &entity.User{
		Username: in.Username,
		Password: pwd,
	}
	err = userModel.CreateUser(u)
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
