package config

import "github.com/zeromicro/go-zero/zrpc"

var C Config

type Config struct {
	zrpc.RpcServerConf
	Mysql struct {
		DataSource      string
		MaxOpenConns    int
		MaxIdleConns    int
		ConnMaxLifetime int
	}
	JWT struct {
		Secret        string
		AccessExpire  int64
		RefreshExpire int64
	}
}
