package model

import (
	"errors"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username         string `db:"username"`
	Password         string `db:"password"`
	Nickname         string `db:"nickname"`
	UserReferralCode string `db:"user_referral_code"`
	ReferralCode     string `db:"referral_code"`
	Wallet           string `db:"wallet"`
	TokenVersion     int    `db:"token_version"`
}

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

func GetByUsername(db *gorm.DB, username string) (*User, error) {
	var user User
	res := db.Table("sys_users").Where("username = ?", username).First(&user)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, ErrUserNotFound
	}
	return &user, nil
}

func GetById(db *gorm.DB, id int) (*User, error) {
	var user User
	res := db.Table("sys_users").Where("id = ?", id).First(&user)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, res.Error
		}
		return nil, ErrUserNotFound
	}
	return &user, nil
}

func CreateUser(db *gorm.DB, user *User) error {
	return db.Table("sys_users").Create(&user).Error
}

func UpdateUser(db *gorm.DB, user *User) error {
	return db.Table("sys_users").Save(&user).Error
}
