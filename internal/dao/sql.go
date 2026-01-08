package dao

import (
	"context"
	"github.com/xxx-newbee/user/internal/config"
	"github.com/xxx-newbee/user/internal/model"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitMysql(c config.Config) {
	dsn := c.Mysql.DataSource
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database: " + err.Error())
	}
	sqlDb, err := db.DB()
	if err != nil {
		panic("failed to get database: " + err.Error())
	}
	sqlDb.SetMaxOpenConns(c.Mysql.MaxOpenConns)
	sqlDb.SetMaxIdleConns(c.Mysql.MaxIdleConns)
	sqlDb.SetConnMaxLifetime(time.Duration(c.Mysql.ConnMaxLifetime) * time.Second)

	if sqlDb.PingContext(context.Background()) != nil {
		panic("failed to ping database: " + err.Error())
	}
	println("✅ MySQL connected successfully")
	db.AutoMigrate(&model.User{})
}

func GetDB() *gorm.DB {
	return db
}

func AutoMigrate(models ...interface{}) error {
	return db.AutoMigrate(models...)
}
