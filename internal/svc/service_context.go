package svc

import (
	"context"
	"image/color"
	"time"

	"github.com/mojocn/base64Captcha"
	"github.com/redis/go-redis/v9"
	"github.com/xxx-newbee/storage"
	"github.com/xxx-newbee/storage/cache"
	"github.com/xxx-newbee/storage/queue"
	"github.com/xxx-newbee/user/internal/config"
	"github.com/xxx-newbee/user/internal/svc/captcha"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config       config.Config
	MemoryQueue  storage.AdapterQueue
	RedisQueue   storage.AdapterQueue
	Database     *gorm.DB
	Cache        storage.AdapterCache
	Captcha      *base64Captcha.Captcha
	CaptchaStore base64Captcha.Store
}

func NewServiceContext(c config.Config) *ServiceContext {
	cacheAdapter := InitRedis(c)
	captchaStore := captcha.NewCaptchaStore(cacheAdapter, c.Captcha.Expire)
	captchaDriver := base64Captcha.NewDriverMath(
		c.Captcha.ImgHeight,
		c.Captcha.ImgWidth,
		c.Captcha.NoiseCount,
		c.Captcha.InterferenceCount,
		&color.RGBA{240, 240, 246, 246},
		base64Captcha.DefaultEmbeddedFonts,
		[]string{"wqy-microhei.ttc"},
	)

	return &ServiceContext{
		Config:       c,
		MemoryQueue:  queue.NewMemoryQueue(c.Queue.Memory.PoolSize),
		RedisQueue:   InitRedisQueue(c),
		Database:     InitDB(c),
		Cache:        cacheAdapter,
		Captcha:      base64Captcha.NewCaptcha(captchaDriver, captchaStore),
		CaptchaStore: captchaStore,
	}
}

func InitDB(c config.Config) *gorm.DB {
	dsn := c.Database.DataSource
	var err error
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database: " + err.Error())
	}
	sqlDb, err := db.DB()
	if err != nil {
		panic("failed to get database: " + err.Error())
	}
	sqlDb.SetMaxOpenConns(c.Database.MaxOpenConns)
	sqlDb.SetMaxIdleConns(c.Database.MaxIdleConns)
	sqlDb.SetConnMaxLifetime(time.Duration(c.Database.ConnMaxLifetime) * time.Second)

	if sqlDb.PingContext(context.Background()) != nil {
		panic("failed to ping database: " + err.Error())
	}
	println("✅ MySQL connected successfully")
	return db
}

func InitRedis(c config.Config) storage.AdapterCache {
	newRedis, err := cache.NewRedis(nil, redis.Options{
		Addr:     c.Cache.Redis.Addr,
		Password: c.Cache.Redis.Password,
		DB:       c.Cache.Redis.DB,
	})
	if err != nil {
		panic("failed to init redis: " + err.Error())
	}

	return newRedis
}

func InitRedisQueue(c config.Config) storage.AdapterQueue {
	client := redis.NewClient(&redis.Options{
		Addr:     c.Queue.Redis.Addr,
		Password: c.Queue.Redis.Password,
		DB:       c.Queue.Redis.DB,
	})
	return queue.NewRedisQueue(client, c.Queue.Redis.Prefix, c.Queue.Redis.MaxRetry)
}
