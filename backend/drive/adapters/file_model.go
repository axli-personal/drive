package adapters

import (
	"github.com/axli-personal/drive/backend/drive/domain"
	"github.com/google/uuid"
	"time"
)

type FileModel struct {
	Id             string `gorm:"size:36;primaryKey"`
	DriveId        string `gorm:"size:36;index:idx_drive_parent_name,unique,priority:1"`
	Parent         string `gorm:"size:36;index:idx_drive_parent_name,unique,priority:2"`
	Name           string `gorm:"size:128;index:idx_drive_parent_name,unique,priority:3"`
	State          string
	LastChange     time.Time
	Size           int
	Hash           string
	DownloadCounts int
}

func NewFileModel(file *domain.File) FileModel {
	return FileModel{
		Id:             file.Id().String(),
		DriveId:        file.DriveId().String(),
		Parent:         file.Parent().String(),
		Name:           file.Name(),
		Hash:           file.Hash(),
		State:          file.State().Value(),
		LastChange:     file.LastChange(),
		Size:           file.Size(),
		DownloadCounts: file.DownloadCounts(),
	}
}

func (model FileModel) TableName() string {
	return "files"
}

func (model FileModel) File() (*domain.File, error) {
	id, err := uuid.Parse(model.Id)
	if err != nil {
		return nil, err
	}

	driveId, err := uuid.Parse(model.DriveId)
	if err != nil {
		return nil, err
	}

	parent, err := domain.CreateParent(model.Parent)
	if err != nil {
		return nil, err
	}

	state, err := domain.CreateState(model.State)
	if err != nil {
		return nil, err
	}

	return domain.NewFileFromRepository(
		id,
		model.Size,
		model.Hash,
		model.DownloadCounts,
		driveId,
		parent,
		model.Name,
		state,
		model.LastChange,
	)
}
