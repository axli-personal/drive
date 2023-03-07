package adapters

import (
	"github.com/axli-personal/drive/internal/drive/domain"
	"github.com/google/uuid"
)

type DriveModel struct {
	Id        string `gorm:"size:36;primaryKey"`
	Owner     string
	UsedBytes int
	PlanName  string
	MaxBytes  int
}

func NewDriveModel(drive *domain.Drive) DriveModel {
	return DriveModel{
		Id:        drive.Id().String(),
		Owner:     drive.Owner(),
		UsedBytes: drive.Usage().Bytes(),
		PlanName:  drive.Plan().Name(),
		MaxBytes:  drive.Plan().MaxBytes(),
	}
}

func (model DriveModel) TableName() string {
	return "drives"
}

func (model DriveModel) Drive() (*domain.Drive, error) {
	id, err := uuid.Parse(model.Id)
	if err != nil {
		return nil, err
	}

	usage, err := domain.NewStorageUsage(model.UsedBytes)
	if err != nil {
		return nil, err
	}

	plan, err := domain.NewStoragePlan(model.PlanName, model.MaxBytes)
	if err != nil {
		return nil, err
	}

	return domain.NewDriveFromRepository(
		id,
		model.Owner,
		usage,
		plan,
	)
}
