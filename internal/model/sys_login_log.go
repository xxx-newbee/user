package model

import (
	"encoding/json"
	"time"

	"github.com/xxx-newbee/storage"
	"gorm.io/gorm"
)

type (
	SysLoginLogModel struct {
		gorm.Model
		Username      string    `json:"username" gorm:"size:128;comment:用户名"`
		Status        string    `json:"status" gorm:"size:4;comment:状态"`
		Ipaddr        string    `json:"ipaddr" gorm:"size:255;comment:ip地址"`
		LoginLocation string    `json:"loginLocation" gorm:"size:255;comment:归属地"`
		Browser       string    `json:"browser" gorm:"size:255;comment:浏览器"`
		Os            string    `json:"os" gorm:"size:255;comment:系统"`
		Platform      string    `json:"platform" gorm:"size:255;comment:固件"`
		LoginTime     time.Time `json:"loginTime" gorm:"comment:登录时间"`
		Remark        string    `json:"remark" gorm:"size:255;comment:备注"`
		Msg           string    `json:"msg" gorm:"size:255;comment:信息"`
	}

	SysLoginLog interface {
		SaveLoginLog(storage.Messager) error
	}

	defaultLoginLogModel struct {
		db    *gorm.DB
		table string
	}
)

func (SysLoginLogModel) TableName() string { return "sys_login_log" }

func NewSysLoginLog(db *gorm.DB) SysLoginLog {
	return &defaultLoginLogModel{
		db:    db,
		table: "sys_login_log",
	}
}

func (o *defaultLoginLogModel) SaveLoginLog(msg storage.Messager) error {
	rb, err := json.Marshal(msg.GetValues())
	if err != nil {
		return err
	}
	var ll SysLoginLog
	if err = json.Unmarshal(rb, &ll); err != nil {
		return err
	}

	if err = o.db.Create(&ll).Error; err != nil {
		return err
	}

	return nil
}
