package model

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/xxx-newbee/storage"
	"gorm.io/gorm"
)

type (
	SysLoginLogModel struct {
		gorm.Model
		UserId        uint64    `json:"userId" gorm:"size:128;comment:用户id"`
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
		DeleteLoginLog(id uint64) error
		GetUserLoginLogs(id uint64, page int) ([]SysLoginLogModel, error)
		GetLoginLogById(id uint64) (SysLoginLogModel, error)
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
	var ll SysLoginLogModel
	if err = json.Unmarshal(rb, &ll); err != nil {
		return err
	}

	if err = o.db.Create(&ll).Error; err != nil {
		return err
	}

	return nil
}

func (o *defaultLoginLogModel) DeleteLoginLog(id uint64) error {
	return o.db.Delete(&SysLoginLogModel{Model: gorm.Model{ID: uint(id)}}).Error
}

func (o *defaultLoginLogModel) GetUserLoginLogs(id uint64, page int) ([]SysLoginLogModel, error) {
	if id == 0 {
		return []SysLoginLogModel{}, nil
	}
	var ll []SysLoginLogModel

	if err := o.db.Table(o.table).Where("user_id = ?", id).Order("created_at desc").Offset((page - 1) * 10).Limit(10).Find(&ll).Error; err != nil {
		return []SysLoginLogModel{}, err
	}
	return ll, nil
}

func (o *defaultLoginLogModel) GetLoginLogById(id uint64) (SysLoginLogModel, error) {
	if id == 0 {
		return SysLoginLogModel{}, nil
	}

	var ll SysLoginLogModel
	if err := o.db.Table(o.table).Where("id = ?", id).First(&ll).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return SysLoginLogModel{}, nil
		}
		return SysLoginLogModel{}, err
	}
	return ll, nil
}
