package adapters

import (
	"context"
	"github.com/axli-personal/drive/backend/user/domain"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(connectionString string) (UserRepository, error) {
	db, err := gorm.Open(mysql.Open(connectionString))
	if err != nil {
		return UserRepository{}, err
	}

	err = db.AutoMigrate(&UserModel{})
	if err != nil {
		return UserRepository{}, err
	}

	return UserRepository{db: db}, nil
}

func (repo UserRepository) SaveUser(ctx context.Context, user *domain.User) error {
	model := NewUserModel(user)

	return repo.db.Create(&model).Error
}

func (repo UserRepository) GetUser(ctx context.Context, account domain.Account) (*domain.User, error) {
	model := UserModel{}

	err := repo.db.Take(&model, "account = ?", account.String()).Error
	if err != nil {
		return nil, err
	}

	return model.User()
}
