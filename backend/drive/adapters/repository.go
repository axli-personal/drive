package adapters

import (
	"github.com/axli-personal/drive/backend/drive/repository"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type GormRepository struct {
	db *gorm.DB
}

func NewMysqlRepository(dsn string) (repository.Repository, error) {
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&DriveModel{}, &FolderModel{}, &FileModel{})
	if err != nil {
		return nil, err
	}

	return GormRepository{db: db}, nil
}

func (repo GormRepository) Transaction(fn func(repo repository.Repository) error) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		return fn(GormRepository{db: tx})
	})
}

func (repo GormRepository) GetDriveRepo() repository.DriveRepository {
	return GormDriveRepository{db: repo.db}
}

func (repo GormRepository) GetFolderRepo() repository.FolderRepository {
	return GormFolderRepository{db: repo.db}
}

func (repo GormRepository) GetFileRepo() repository.FileRepository {
	return GormFileRepository{db: repo.db}
}
