package model

import (
	"errors"

	"gorm.io/gorm"
)

type (
	User struct {
		gorm.Model
		Username         string `db:"username"`
		Password         string `db:"password"`
		Nickname         string `db:"nickname"`
		UserReferralCode string `db:"user_referral_code"`
		ReferralCode     string `db:"referral_code"`
		Wallet           string `db:"wallet"`
		TokenVersion     int    `db:"token_version"`
	}

	UserModel interface {
		GetByUsername(username string) (*User, error)
		GetById(id uint) (*User, error)
		Create(user *User) error
		Update(user *User) error
	}
)

func (u *User) TableName() string {
	return "sys_users"
}

var (
	ErrUserAlreadyExist            = errors.New("user already exists")
	ErrPasswordNecessary           = errors.New("password is necessary")
	ErrUsernameOrPasswordEmpty     = errors.New("username or password cannot be empty")
	ErrUserCreateFailed            = errors.New("failed to create user")
	ErrUsernameOrPasswordIncorrect = errors.New("username or password is incorrect")
	ErrUpdateUserFailed            = errors.New("failed to update user")
	ErrPasswordEmpty               = errors.New("password is empty")
	ErrUserNotFound                = errors.New("user not found")
	ErrTokenExpired                = errors.New("token is expired")
	ErrOldPasswordIncorrect        = errors.New("old password is incorrect")
	ErrChangePasswordFailed        = errors.New("failed to change password")
	ErrGenerateReferralCode        = errors.New("failed to generate referral code")
	ErrCaptchaIncorrect            = errors.New("captcha is incorrect")
)

type defaultUserModel struct {
	db    *gorm.DB
	table string
}

func NewDefaultUser(db *gorm.DB) UserModel {
	return &defaultUserModel{
		db:    db,
		table: "sys_users",
	}
}

func (o *defaultUserModel) GetByUsername(username string) (*User, error) {
	var user User
	res := o.db.Table(o.table).Where("username = ?", username).First(&user)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, res.Error
	}
	return &user, nil
}

func (o *defaultUserModel) GetById(id uint) (*User, error) {
	var user User
	res := o.db.Table(o.table).Where("id = ?", id).First(&user)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, res.Error
	}
	return &user, nil
}

func (o *defaultUserModel) Create(user *User) error {
	return o.db.Table(o.table).Create(&user).Error
}

func (o *defaultUserModel) Update(user *User) error {
	return o.db.Table(o.table).Save(&user).Error
}
