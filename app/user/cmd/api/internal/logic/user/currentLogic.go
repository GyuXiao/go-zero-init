package user

import (
	"context"
	"go-zero-init/app/user/cmd/rpc/client/user"
	"strings"

	"go-zero-init/app/user/cmd/api/internal/svc"
	"go-zero-init/app/user/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CurrentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCurrentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CurrentLogic {
	return &CurrentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CurrentLogic) Current(req *types.CurrentUserReq) (resp *types.CurrentUserResp, err error) {
	// 获取 JWT token
	token := strings.Split(req.Authorization, " ")[1]

	// 调用 rpc 模块的 current
	currentResp, err := l.svcCtx.UserRpc.CurrentUser(l.ctx, &user.CurrentUserReq{AuthToken: token})
	if err != nil {
		return nil, err
	}

	// 返回响应参数
	return &types.CurrentUserResp{
		Id:          currentResp.Id,
		Username:    currentResp.Username,
		AvatarUrl:   currentResp.AvatarUrl,
		UserRole:    uint8(currentResp.UserRole),
		Token:       currentResp.Token,
		TokenExpire: currentResp.TokenExpire,
	}, nil
}
