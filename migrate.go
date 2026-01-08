package main

import (
	"flag"
	"github.com/xxx-newbee/go-micro/user/internal/config"
	"github.com/xxx-newbee/go-micro/user/internal/dao"
	"github.com/xxx-newbee/go-micro/user/internal/model"
	"log"

	"github.com/zeromicro/go-zero/core/conf"
)

var configFile = flag.String("f", "etc/user.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	dao.InitMysql(c)
	if err := dao.AutoMigrate(&model.User{}); err != nil {
		log.Fatal(err)
	}
}
