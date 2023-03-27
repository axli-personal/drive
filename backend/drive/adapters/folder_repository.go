package adapters

import (
	"context"
	"errors"
	"github.com/axli-personal/drive/backend/drive/domain"
	"github.com/axli-personal/drive/backend/drive/repository"
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
	return saveFolder(ctx, repo.db, folder)
}

func (repo GormFolderRepository) GetFolder(ctx context.Context, id uuid.UUID) (*domain.Folder, error) {
	return getFolder(ctx, repo.db, id)
}

func (repo GormFolderRepository) FindFolder(ctx context.Context, options repository.FindFolderOptions) ([]*domain.Folder, error) {
	return findFolder(ctx, repo.db, options)
}

func (repo GormFolderRepository) UpdateFolder(ctx context.Context, folder *domain.Folder, options repository.UpdateFolderOptions) error {
	return updateFolder(ctx, repo.db, folder, options)
}

func saveFolder(ctx context.Context, tx *gorm.DB, folder *domain.Folder) error {
	folderModel := NewFolderModel(folder)

	return tx.Create(&folderModel).Error
}

func getFolder(ctx context.Context, tx *gorm.DB, id uuid.UUID) (*domain.Folder, error) {
	folderModel := FolderModel{}

	err := tx.Take(&folderModel, "id = ?", id.String()).Error
	if err != nil {
		return nil, err
	}

	return folderModel.Folder()
}

func findFolder(ctx context.Context, tx *gorm.DB, options repository.FindFolderOptions) ([]*domain.Folder, error) {
	var folderModels []FolderModel

	err := tx.Where("parent = ?", options.Parent.String()).Find(&folderModels).Error
	if err != nil {
		return nil, err
	}

	var folders []*domain.Folder
	for _, model := range folderModels {
		folder, err := model.Folder()
		if err != nil {
			return nil, err
		}
		folders = append(folders, folder)
	}

	return folders, err
}

func updateFolder(ctx context.Context, tx *gorm.DB, folder *domain.Folder, options repository.UpdateFolderOptions) error {
	folderModel := NewFolderModel(folder)

	err := tx.Save(&folderModel).Error
	if err != nil {
		return err
	}

	if options.UpdateChildrenState {
		parent, err := domain.CreateFolderParent(folder.Id())
		if err != nil {
			return err
		}

		subFiles, err := findFile(ctx, tx, repository.FindFileOptions{Parent: parent})
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		for i := 0; i < len(subFiles); i++ {
			subFiles[i].SetState(folder.State())
			err = updateFile(ctx, tx, subFiles[i], repository.UpdateFileOptions{})
			if err != nil {
				return err
			}
		}

		subFolders, err := findFolder(ctx, tx, repository.FindFolderOptions{Parent: parent})
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		for i := 0; i < len(subFolders); i++ {
			subFolders[i].SetState(folder.State())
			err = updateFolder(ctx, tx, subFolders[i], repository.UpdateFolderOptions{UpdateChildrenState: true})
			if err != nil {
				return err
			}
		}
	}

	return nil
}
