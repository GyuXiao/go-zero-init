package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"go-zero-init/app/user/cmd/api/internal/config"
	"go-zero-init/app/user/cmd/rpc/client/user"
)

type ServiceContext struct {
	Config      config.Config
	RedisClient *redis.Redis
	UserRpc     user.UserZrpcClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		RedisClient: redis.MustNewRedis(redis.RedisConf{
			Host: c.Redis.Host,
			Type: c.Redis.Type,
			Pass: c.Redis.Pass,
		}),
		UserRpc: user.NewUserZrpcClient(zrpc.MustNewClient(c.UserRpcConf)),
	}
}
