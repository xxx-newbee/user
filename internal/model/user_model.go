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
	ErrPasswordNeccessary          = errors.New("password is necessary")
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
)
