package adapters

import (
	"github.com/axli-personal/drive/backend/drive/domain"
	"github.com/google/uuid"
	"time"
)

type FolderModel struct {
	Id         string `gorm:"size:36;primaryKey"`
	DriveId    string `gorm:"size:36;index:idx_drive_parent_name,unique,priority:1"`
	Parent     string `gorm:"size:36;index:idx_drive_parent_name,unique,priority:2"`
	Name       string `gorm:"size:128;index:idx_drive_parent_name,unique,priority:3"`
	State      string
	LastChange time.Time
}

func NewFolderModel(folder *domain.Folder) FolderModel {
	return FolderModel{
		Id:         folder.Id().String(),
		DriveId:    folder.DriveId().String(),
		Parent:     folder.Parent().String(),
		Name:       folder.Name(),
		State:      folder.State().Value(),
		LastChange: folder.LastChange(),
	}
}

func (model FolderModel) TableName() string {
	return "folders"
}

func (model FolderModel) Folder() (*domain.Folder, error) {
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

	return domain.NewFolderFromRepository(
		id,
		driveId,
		parent,
		model.Name,
		state,
		model.LastChange,
	)
}
