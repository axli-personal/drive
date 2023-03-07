package adapters

import (
	"context"
	"github.com/axli-personal/drive/internal/drive/domain"
	"github.com/axli-personal/drive/internal/drive/repository"
	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type GormFolderRepository struct {
	db *gorm.DB
}

func NewMysqlFolderRepository(connectionString string) (repository.FolderRepository, error) {
	db, err := gorm.Open(mysql.Open(connectionString))
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&FolderModel{})
	if err != nil {
		return nil, err
	}

	return GormFolderRepository{db: db}, nil
}

func (repo GormFolderRepository) SaveFolder(ctx context.Context, folder *domain.Folder) error {
	model := NewFolderModel(folder)

	return repo.db.Create(&model).Error
}

func (repo GormFolderRepository) GetFolder(ctx context.Context, id uuid.UUID) (*domain.Folder, error) {
	model := FolderModel{}

	err := repo.db.Take(&model, "id = ?", id.String()).Error
	if err != nil {
		return nil, err
	}

	return model.Folder()
}

func (repo GormFolderRepository) FindFolder(ctx context.Context, options repository.FindFolderOptions) ([]*domain.Folder, error) {
	var models []FolderModel

	err := repo.db.Where("parent = ?", options.Parent.String()).Find(&models).Error
	if err != nil {
		return nil, err
	}

	var folders []*domain.Folder
	for _, model := range models {
		folder, err := model.Folder()
		if err != nil {
			return nil, err
		}
		folders = append(folders, folder)
	}

	return folders, err
}

func (repo GormFolderRepository) UpdateFolder(ctx context.Context, foler *domain.Folder, options repository.UpdateFolderOptions) error {
	model := NewFolderModel(foler)

	var err error
	if options.MustInSameState {
		err = repo.db.Where("state = ?", model.State).Save(&model).Error
	} else if !options.MustInState.IsZero() {
		err = repo.db.Where("state = ?", options.MustInState.Value()).Save(&model).Error
	} else {
		err = repo.db.Save(&model).Error
	}

	return err
}
