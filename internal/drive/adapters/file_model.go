package adapters

import (
	"github.com/axli-personal/drive/internal/drive/domain"
	"github.com/google/uuid"
	"time"
)

type FileModel struct {
	Id             string `gorm:"size:36;primaryKey"`
	DriveId        string
	Parent         string
	Name           string
	State          string
	LastChange     time.Time
	Endpoint       string
	Bytes          int
	DownloadCounts int
}

func NewFileModel(file *domain.File) FileModel {
	return FileModel{
		Id:             file.Id().String(),
		DriveId:        file.DriveId().String(),
		Parent:         file.Parent().String(),
		Name:           file.Name(),
		State:          file.State().Value(),
		LastChange:     file.LastChange(),
		Endpoint:       file.Endpoint(),
		Bytes:          file.Bytes(),
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
		model.Endpoint,
		model.Bytes,
		model.DownloadCounts,
		driveId,
		parent,
		model.Name,
		state,
		model.LastChange,
	)
}
