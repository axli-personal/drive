package adapters

import (
	"context"
	"github.com/axli-personal/drive/internal/drive/domain"
	"github.com/axli-personal/drive/internal/drive/repository"
	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type GormFileRepository struct {
	db *gorm.DB
}

func NewMysqlFileRepository(connectString string) (repository.FileRepository, error) {
	db, err := gorm.Open(mysql.Open(connectString))
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&FileModel{})
	if err != nil {
		return nil, err
	}

	return GormFileRepository{db: db}, nil
}

func (repo GormFileRepository) SaveFile(ctx context.Context, file *domain.File) error {
	model := NewFileModel(file)

	return repo.db.Create(&model).Error
}

func (repo GormFileRepository) GetFile(ctx context.Context, id uuid.UUID) (*domain.File, error) {
	model := FileModel{}

	err := repo.db.Take(&model, "id = ?", id.String()).Error
	if err != nil {
		return nil, err
	}

	return model.File()
}

func (repo GormFileRepository) FindFile(ctx context.Context, options repository.FindFileOptions) ([]*domain.File, error) {
	var models []FileModel

	err := repo.db.Where("parent = ?", options.Parent.String()).Find(&models).Error
	if err != nil {
		return nil, err
	}

	var files []*domain.File
	for _, model := range models {
		file, err := model.File()
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	return files, err
}

func (repo GormFileRepository) UpdateFile(ctx context.Context, file *domain.File, options repository.UpdateFileOptions) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		fileModel := NewFileModel(file)

		var err error
		if options.MustInSameState {
			err = tx.Where("state = ?", fileModel.State).Save(&fileModel).Error
		} else if !options.MustInState.IsZero() {
			err = tx.Where("state = ?", options.MustInState.Value()).Save(&fileModel).Error
		} else {
			err = tx.Save(&fileModel).Error
		}
		if err != nil {
			return err
		}

		if options.IncreaseStorageUsage {
			err = tx.Model(&DriveModel{}).Where("id = ?", fileModel.DriveId).Update("used_bytes", gorm.Expr("used_bytes + ?", fileModel.Bytes)).Error
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (repo GormFileRepository) DeleteFile(ctx context.Context, file *domain.File) error {
	model := NewFileModel(file)

	return repo.db.Delete(&model).Error
}
