package adapters

import (
	"context"
	"errors"
	"github.com/axli-personal/drive/internal/drive/domain"
	"github.com/axli-personal/drive/internal/drive/repository"
	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	ErrMultipleDrive = errors.New("multiple drive")
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
				return ErrMultipleDrive
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
		return nil, err
	}

	return model.Drive()
}

func (repo GormDriveRepository) UpdateDrive(ctx context.Context, drive *domain.Drive) error {
	model := NewDriveModel(drive)

	return repo.db.Save(&model).Error
}
