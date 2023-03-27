package adapters

import (
	"context"
	"github.com/axli-personal/drive/backend/storage/domain"
	"github.com/axli-personal/drive/backend/storage/repository"
	"io"
	"os"
	"path"
)

type DiskObjectRepository struct {
	dirPath string
}

func NewDiskObjectRepository(directoryPath string) (repository.ObjectRepository, error) {
	err := os.MkdirAll(directoryPath, 0700)
	if err != nil {
		return nil, err
	}

	return DiskObjectRepository{
		dirPath: directoryPath,
	}, nil
}

func (repo DiskObjectRepository) SaveObject(ctx context.Context, object *domain.Object) error {
	objectPath := path.Join(repo.dirPath, object.Hash())

	_, err := os.Stat(objectPath)
	if os.IsExist(err) {
		return nil
	}

	file, err := os.Create(objectPath)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = io.Copy(file, object)
	if err != nil {
		return err
	}

	return nil
}

func (repo DiskObjectRepository) GetObject(ctx context.Context, hash string) (*domain.Object, error) {
	file, err := os.Open(path.Join(repo.dirPath, hash))
	if err != nil {
		return nil, err
	}

	return domain.NewObject(file)
}
