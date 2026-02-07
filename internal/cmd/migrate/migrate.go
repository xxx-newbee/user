package migrate

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/xxx-newbee/user/internal/config"
	"github.com/xxx-newbee/user/internal/model"
	"github.com/xxx-newbee/user/internal/svc"
	"github.com/zeromicro/go-zero/core/conf"
)

var (
	configYaml string
	StartCmd   = &cobra.Command{
		Use:     "migrate",
		Short:   "migrate command",
		Example: "./user-srv migrate -c etc/user.yaml",
		PreRun: func(cmd *cobra.Command, args []string) {
			setup()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
)

func init() {
	StartCmd.PersistentFlags().StringVarP(&configYaml, "config", "c", "etc/user.yaml", "config file")
}

func setup() {
	conf.MustLoad(configYaml, &config.C)
}

func run() error {
	db := svc.InitDB(config.C)
	if err := db.AutoMigrate(&model.User{}, &model.SysLoginLog{}); err != nil {
		return err
	}
	log.Println("migrate success")
	return nil
}
