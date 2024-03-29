package main

import (
	"flag"
	"fmt"

	"github.com/finishy1995/effibot-core/ai/internal/config"
	"github.com/finishy1995/effibot-core/ai/internal/server"
	"github.com/finishy1995/effibot-core/ai/internal/svc"
	"github.com/finishy1995/effibot-core/ai/pb/ai"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/ai.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		ai.RegisterAiServer(grpcServer, server.NewAiServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
