package migrate

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/xxx-newbee/user/internal/config"
	"github.com/xxx-newbee/user/internal/dao"
	"github.com/xxx-newbee/user/internal/model"
	"github.com/zeromicro/go-zero/core/conf"
)

var (
	configYaml string
	StartCmd   = &cobra.Command{
		Use:     "migrate",
		Short:   "migrate command",
		Example: "./user-srv migrate -c etc/user.yaml",
		Run: func(cmd *cobra.Command, args []string) {
			run()
		},
	}
)

func init() {
	StartCmd.PersistentFlags().StringVarP(&configYaml, "config", "c", "etc/user.yaml", "config file")
}

func run() {
	var c config.Config
	conf.MustLoad(configYaml, &c)
	dao.InitMysql(c)
	if err := dao.AutoMigrate(&model.User{}); err != nil {
		log.Fatal(err)
	}
	log.Println("migrate success")
}
