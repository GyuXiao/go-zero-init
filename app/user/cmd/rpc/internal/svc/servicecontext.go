package svc

import (
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"go-zero-init/app/user/cmd/rpc/internal/config"
	"go-zero-init/app/user/models/do"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config      config.Config
	DBEngin     *do.Query
	RedisClient *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := gorm.Open(mysql.Open(c.DB.DataSource), &gorm.Config{})
	if err != nil {
		panic("failed to connect database, error=" + err.Error())
	}
	logc.Info(context.Background(), "connect MySQL database success")

	return &ServiceContext{
		Config:  c,
		DBEngin: do.Use(db),
		RedisClient: redis.MustNewRedis(redis.RedisConf{
			Host: c.Redis.Host,
			Type: c.Redis.Type,
			Pass: c.Redis.Pass,
		}),
	}
}
