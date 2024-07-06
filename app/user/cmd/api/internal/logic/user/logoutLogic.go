package user

import (
	"context"
	"go-zero-init/app/user/cmd/rpc/client/user"
	"strings"

	"go-zero-init/app/user/cmd/api/internal/svc"
	"go-zero-init/app/user/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogoutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LogoutLogic) Logout(req *types.LogoutReq) (resp *types.LogoutResp, err error) {
	// 获取 JWT token
	token := strings.Split(req.Authorization, " ")[1]

	// 调用 rpc 模块的 logout
	logoutResp, err := l.svcCtx.UserRpc.Logout(l.ctx, &user.LogoutReq{AuthToken: token})
	if err != nil {
		return nil, err
	}

	// 返回响应参数
	return &types.LogoutResp{IsLogouted: logoutResp.IsLogouted}, nil
}
