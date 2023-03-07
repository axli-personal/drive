package adapters

import (
	"github.com/axli-personal/drive/internal/user/domain"
)

type UserModel struct {
	Account      string `gorm:"size:30;primaryKey"`
	Username     string `gorm:"size:64"`
	Password     string `gorm:"size:30"`
	Introduction string `gorm:"size:256"`
}

func (model UserModel) TableName() string {
	return "users"
}

func NewUserModel(user *domain.User) UserModel {
	return UserModel{
		Account:      user.Account().String(),
		Username:     user.Username(),
		Password:     user.Password().String(),
		Introduction: user.Introduction(),
	}
}

func (model UserModel) User() (*domain.User, error) {
	account, err := domain.NewAccount(model.Account)
	if err != nil {
		return nil, err
	}

	password, err := domain.NewPassword(model.Password)
	if err != nil {
		return nil, err
	}

	return domain.NewUserFromRepository(
		account,
		password,
		model.Username,
		model.Introduction,
	)
}
