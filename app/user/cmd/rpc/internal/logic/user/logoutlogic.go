package userlogic

import (
	"context"
	"go-zero-init/app/user/models/dao"

	"go-zero-init/app/user/cmd/rpc/internal/svc"
	"go-zero-init/app/user/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogoutLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LogoutLogic) Logout(in *pb.LogoutReq) (*pb.LogoutResp, error) {
	tokenLogic := dao.NewDefaultTokenModel(l.svcCtx.RedisClient)
	err := tokenLogic.DeleteToken(in.AuthToken)
	if err != nil {
		return nil, err
	}

	return &pb.LogoutResp{IsLogouted: true}, nil
}
