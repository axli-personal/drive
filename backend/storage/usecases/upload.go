package usecases

import (
	"context"
	"github.com/axli-personal/drive/backend/common/decorator"
	"github.com/axli-personal/drive/backend/pkg/types"
	"github.com/axli-personal/drive/backend/storage/domain"
	"github.com/axli-personal/drive/backend/storage/remote"
	"github.com/axli-personal/drive/backend/storage/repository"
	"github.com/sirupsen/logrus"
	"io"
)

type (
	UploadArgs struct {
		SessionId  string
		FileParent string
		FileName   string
		Data       io.ReadSeeker
	}

	UploadResult struct {
	}
)

type uploadHandler struct {
	driveService remote.DriveService
	objectRepo   repository.ObjectRepository
	capacityRepo repository.CapacityRepository
}

func (handler uploadHandler) Handle(ctx context.Context, args UploadArgs) (UploadResult, error) {
	object, err := domain.NewObject(args.Data)
	if err != nil {
		return UploadResult{}, err
	}

	startResponse, err := handler.driveService.StartUpload(
		types.StartUploadRequest{
			SessionId:  args.SessionId,
			FileParent: args.FileParent,
			FileName:   args.FileName,
			FileHash:   object.Hash(),
			FileSize:   object.TotalBytes(),
		},
	)
	if err != nil {
		return UploadResult{}, err
	}

	err = handler.objectRepo.SaveObject(ctx, object)
	if err != nil {
		return UploadResult{}, err
	}

	_, err = handler.driveService.FinishUpload(
		types.FinishUploadRequest{
			FileId: startResponse.FileId,
		},
	)
	if err != nil {
		return UploadResult{}, err
	}

	return UploadResult{}, nil
}

type UploadHandler decorator.Handler[UploadArgs, UploadResult]

func NewUploadHandler(
	driveService remote.DriveService,
	objectRepo repository.ObjectRepository,
	logger *logrus.Entry,
) UploadHandler {
	return decorator.WithLogging[UploadArgs, UploadResult](
		uploadHandler{
			driveService: driveService,
			objectRepo:   objectRepo,
		},
		logger,
	)
}
