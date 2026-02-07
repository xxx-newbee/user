package config

import (
	"github.com/zeromicro/go-zero/zrpc"
)

var C Config

type Config struct {
	zrpc.RpcServerConf
	Database struct {
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
	Cache struct {
		Redis  RedisConf
		Memory string
	}
	Queue struct {
		Redis  RedisQueueConf
		Memory MemoryQueueConf
	}
	Captcha CaptchaConf
}

type RedisConf struct {
	Addr     string
	Password string
	DB       int
}

type RedisQueueConf struct {
	Addr     string
	Password string
	DB       int
	Prefix   string
	MaxRetry int
}

type MemoryQueueConf struct {
	PoolSize int
}

type CaptchaConf struct {
	Expire            int `json:",default=600"`
	FontSize          int `json:",default=30"`
	ImgWidth          int `json:",default=80"`
	ImgHeight         int `json:",default=30"`
	NoiseCount        int `json:",default=1"`
	InterferenceCount int `json:",default=0"`
	MathDifficulty    int `json:",default=10"`
}
