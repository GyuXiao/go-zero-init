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

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	// 校验参数
	if req.Username == constant.BlankString || req.Password == constant.BlankString || len(req.Username) < constant.UsernameMinLen || len(req.Password) < constant.PasswordMinLen {
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "用户名或密码错误")
	}
	_, err = regexp.MatchString(constant.PatternStr, req.Username)
	if err != nil {
		return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "用户名称包含非法字符")
	}

	// 调用 rpc 模块的 login
	loginResp, err := l.svcCtx.UserRpc.Login(l.ctx, &user.LoginReq{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	// 返回登陆响应参数
	return &types.LoginResp{
		Id:          loginResp.Id,
		Username:    loginResp.Username,
		AvatarUrl:   loginResp.AvatarUrl,
		UserRole:    uint8(loginResp.UserRole),
		Token:       loginResp.Token,
		TokenExpire: loginResp.TokenExpire,
	}, nil
}
