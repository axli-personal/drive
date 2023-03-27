package adapters

import (
	"context"
	"github.com/axli-personal/drive/backend/drive/domain"
	"github.com/axli-personal/drive/backend/drive/repository"
	"github.com/axli-personal/drive/backend/pkg/errors"
	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type GormDriveRepository struct {
	db *gorm.DB
}

func NewMysqlDriveRepository(connectString string) (repository.DriveRepository, error) {
	db, err := gorm.Open(mysql.Open(connectString))
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&DriveModel{})
	if err != nil {
		return nil, err
	}

	return GormDriveRepository{db: db}, nil
}

func (repo GormDriveRepository) CreateDrive(ctx context.Context, drive *domain.Drive, options repository.CreateDriveOptions) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		model := NewDriveModel(drive)

		if options.OnlyOneDrive {
			err := tx.Where("owner = ?", model.Owner).Take(&DriveModel{}).Error
			if err == nil {
				return errors.New(repository.ErrCodeRepository, "multiple drive", nil)
			}
			if err != gorm.ErrRecordNotFound {
				return err
			}
		}

		return repo.db.Create(&model).Error
	})
}

func (repo GormDriveRepository) GetDrive(ctx context.Context, id uuid.UUID) (*domain.Drive, error) {
	model := DriveModel{}

	err := repo.db.Take(&model, "id = ?", id.String()).Error
	if err != nil {
		return nil, err
	}

	return model.Drive()
}

func (repo GormDriveRepository) GetDriveByOwner(ctx context.Context, owner string) (*domain.Drive, error) {
	model := DriveModel{}

	err := repo.db.Take(&model, "owner = ?", owner).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(repository.ErrCodeNotFound, "drive not found", err)
		}
		return nil, errors.New(repository.ErrCodeRepository, "fail to get drive by owner", err)
	}

	return model.Drive()
}

func (repo GormDriveRepository) UpdateDrive(ctx context.Context, drive *domain.Drive) error {
	model := NewDriveModel(drive)

	return repo.db.Save(&model).Error
}
