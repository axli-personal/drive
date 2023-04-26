package adapters

import (
	"context"
	"github.com/axli-personal/drive/backend/drive/domain"
	"github.com/axli-personal/drive/backend/drive/repository"
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
	fileModel := NewFileModel(file)

	return repo.db.WithContext(ctx).Create(&fileModel).Error
}

func (repo GormFileRepository) GetFile(ctx context.Context, id uuid.UUID) (*domain.File, error) {
	return getFile(ctx, repo.db, id)
}

func (repo GormFileRepository) FindFile(ctx context.Context, options repository.FindFileOptions) ([]*domain.File, error) {
	tx := repo.db.WithContext(ctx)

	if options.DriveId != uuid.Nil {
		tx = tx.Where("drive_id = ?", options.DriveId.String())
	}

	if !options.Parent.IsZero() {
		tx = tx.Where("parent = ?", options.Parent.String())
	}

	if options.Name != "" {
		tx = tx.Where("name = ?", options.Name)
	}

	if len(options.States) > 0 {
		var statesCond []string
		for _, state := range options.States {
			statesCond = append(statesCond, state.Value())
		}
		tx = tx.Where("state IN ?", statesCond)
	}

	var models []FileModel
	err := tx.Find(&models).Error
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
	return updateFile(ctx, repo.db, file, options)
}

func (repo GormFileRepository) DeleteFile(ctx context.Context, file *domain.File) error {
	return deleteFile(ctx, repo.db, file)
}

func getFile(ctx context.Context, tx *gorm.DB, id uuid.UUID) (*domain.File, error) {
	fileModel := FileModel{}

	err := tx.Take(&fileModel, "id = ?", id.String()).Error
	if err != nil {
		return nil, err
	}

	return fileModel.File()
}

func updateFile(ctx context.Context, tx *gorm.DB, file *domain.File, options repository.UpdateFileOptions) error {
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
		err = tx.Model(&DriveModel{}).Where("id = ?", fileModel.DriveId).Update("used_bytes", gorm.Expr("used_bytes + ?", fileModel.Size)).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func deleteFile(ctx context.Context, tx *gorm.DB, file *domain.File) error {
	model := NewFileModel(file)

	return tx.Delete(&model).Error
}
