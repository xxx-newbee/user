package api

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xxx-newbee/user/internal/config"
	"github.com/xxx-newbee/user/internal/logic"
	"github.com/xxx-newbee/user/internal/model"
	"github.com/xxx-newbee/user/internal/server"
	"github.com/xxx-newbee/user/internal/svc"
	"github.com/xxx-newbee/user/user"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	configYaml string
	StartCmd   = &cobra.Command{
		Use:          "service",
		Short:        "start service",
		Example:      "go run main.go service -c config.yaml",
		SilenceUsage: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			setup()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
)

func init() {
	StartCmd.PersistentFlags().StringVarP(&configYaml, "config", "c", "etc/user.yaml", "the config file")
}

func setup() {
	conf.MustLoad(configYaml, &config.C)
}

func run() error {
	sctx := svc.NewServiceContext(config.C)
	// 注册登录日志消费者
	sctx.MemoryQueue.Register(model.SysLoginLog{}.TableName(), logic.NewLoginLogic(context.Background(), sctx).SaveLoginLog)
	go sctx.MemoryQueue.Run()

	s := zrpc.MustNewServer(config.C.RpcServerConf, func(grpcServer *grpc.Server) {
		user.RegisterUserServer(grpcServer, server.NewUserServer(sctx))

		if config.C.Mode == service.DevMode || config.C.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer func() {
		sctx.MemoryQueue.Shutdown()
		s.Stop()
	}()

	fmt.Printf("Starting rpc server at %s...\n", config.C.ListenOn)
	s.Start()
	return nil
}
