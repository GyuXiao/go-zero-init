package main

import (
	"flag"
	"fmt"
	"go-zero-init/common/constant"
	"go-zero-init/common/result"
	"net/http"

	"go-zero-init/app/user/cmd/api/internal/config"
	"go-zero-init/app/user/cmd/api/internal/handler"
	"go-zero-init/app/user/cmd/api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/user.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf,
		rest.WithCustomCors(nil, func(w http.ResponseWriter) {}, constant.AllOrigins),
		rest.WithUnauthorizedCallback(result.JwtUnauthorizedResult))
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
