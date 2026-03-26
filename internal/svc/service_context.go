package svc

import (
	"context"
	"image/color"
	"time"

	"github.com/mojocn/base64Captcha"
	"github.com/redis/go-redis/v9"
	"github.com/xxx-newbee/storage"
	"github.com/xxx-newbee/storage/cache"
	"github.com/xxx-newbee/storage/locker"
	"github.com/xxx-newbee/storage/queue"
	"github.com/xxx-newbee/user/internal/config"
	"github.com/xxx-newbee/user/internal/model"
	"github.com/xxx-newbee/user/internal/svc/captcha"
	"github.com/xxx-newbee/user/internal/svc/mail"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config      config.Config
	Locker      storage.AdapterLocker
	MemoryQueue storage.AdapterQueue
	RedisQueue  storage.AdapterQueue

	// verifies
	Cache        storage.AdapterCache
	Captcha      *base64Captcha.Captcha
	CaptchaStore base64Captcha.Store
	MailStore    mail.MailStore

	// models
	UserModel   model.UserModel
	SysLoginLog model.SysLoginLog
}

func NewServiceContext(c config.Config) *ServiceContext {
	db := InitDB(c)
	rdb := InitRedis(c)
	cacheAdapter := cache.NewRedis(rdb, nil)
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
		RedisQueue:   queue.NewRedisQueue(rdb, c.Queue.Redis.Prefix, c.Queue.Redis.MaxRetry),
		Cache:        cacheAdapter,
		Locker:       locker.NewRedisLocker(rdb),
		Captcha:      base64Captcha.NewCaptcha(captchaDriver, captchaStore),
		CaptchaStore: captchaStore,
		MailStore:    mail.NewMailVerifyStore(c, cacheAdapter, c.SMTP.Expire),
		UserModel:    model.NewDefaultUser(db),
		SysLoginLog:  model.NewSysLoginLog(db),
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
	println("Database connected！✅")
	return db
}

func InitRedis(c config.Config) *redis.Client {

	newRedis := redis.NewClient(&redis.Options{
		Addr:     c.Cache.Redis.Addr,
		Password: c.Cache.Redis.Password,
		DB:       c.Cache.Redis.DB,
	})
	if err := newRedis.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}
	return newRedis
}
