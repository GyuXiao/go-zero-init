package user

import (
	"context"
	"go-zero-init/app/user/cmd/rpc/client/user"
	"go-zero-init/common/constant"
	"go-zero-init/common/xerr"
	"regexp"

	"go-zero-init/app/user/cmd/api/internal/svc"
	"go-zero-init/app/user/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
	// 校验参数
	if req.Username == constant.BlankString || req.Password == constant.BlankString || len(req.Username) < constant.UsernameMinLen || len(req.Password) < constant.PasswordMinLen || req.Password != req.ConfirmPassword {
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "用户名或密码错误")
	}
	_, err = regexp.MatchString(constant.PatternStr, req.Username)
	if err != nil {
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "用户名称包含非法字符")
	}

	// 调用 rpc 模块的 register
	registerResp, err := l.svcCtx.UserRpc.Register(l.ctx, &user.RegisterReq{
		Username:        req.Username,
		Password:        req.Password,
		ConfirmPassword: req.ConfirmPassword,
	})
	if err != nil {
		return nil, err
	}

	// 返回响应参数
	return &types.RegisterResp{Username: registerResp.Username}, nil
}
