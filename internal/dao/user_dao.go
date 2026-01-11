package dao

import (
	"errors"

	"github.com/xxx-newbee/user/internal/model"

	"gorm.io/gorm"
)

type UserDao struct {
	db *gorm.DB
}

func NewUserDao() *UserDao {
	return &UserDao{db: GetDB()}
}

func (dao *UserDao) GetByUsername(username string) (*model.User, error) {
	var user model.User
	result := dao.db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

func (dao *UserDao) GetByID(id uint) (*model.User, error) {
	var user model.User
	result := dao.db.First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

func (dao *UserDao) Create(user model.User) error {
	result := dao.db.Create(&user)
	return result.Error
}

func (dao *UserDao) Update(user *model.User) error {
	result := dao.db.Save(&user)
	return result.Error
}
