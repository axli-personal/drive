package usecases

import (
	"context"
	"errors"
	"github.com/axli-personal/drive/internal/common/decorator"
	"github.com/axli-personal/drive/internal/pkg/events"
	"github.com/axli-personal/drive/internal/pkg/types"
	"github.com/axli-personal/drive/internal/storage/domain"
	"github.com/axli-personal/drive/internal/storage/remote"
	"github.com/axli-personal/drive/internal/storage/repository"
	"github.com/sirupsen/logrus"
	"io"
)

var (
	ErrCannotUpload = errors.New("cannot upload")
)

type (
	UploadObjectArgs struct {
		SessionId string
		FileId    string
		Data      io.Reader
	}

	UploadObjectResult struct {
		// TODO: do we need to return the total bytes of the file?
	}
)

type uploadObjectHandler struct {
	endpoint     string
	driveService remote.DriveService
	objectRepo   repository.ObjectRepository
	eventRepo    repository.EventRepository
	capacityRepo repository.CapacityRepository
}

func (handler uploadObjectHandler) Handle(ctx context.Context, args UploadObjectArgs) (UploadObjectResult, error) {
	handler.capacityRepo.DecreaseRequestCapacity(ctx, 1)

	response, err := handler.driveService.CanUpload(
		types.CanUploadRequest{
			SessionId: args.SessionId,
			FileId:    args.FileId,
		},
	)
	if err != nil {
		return UploadObjectResult{}, err
	}
	if !response.CanUpload {
		return UploadObjectResult{}, ErrCannotUpload
	}

	object, err := domain.NewObject(args.FileId, args.Data)
	if err != nil {
		return UploadObjectResult{}, err
	}

	err = handler.objectRepo.SaveObject(ctx, object)
	if err != nil {
		return UploadObjectResult{}, err
	}

	err = handler.eventRepo.PublishFileUploaded(
		ctx,
		events.FileUploaded{
			Endpoint:   handler.endpoint,
			FileId:     object.FileId(),
			TotalBytes: object.TotalBytes(),
		},
	)
	if err != nil {
		return UploadObjectResult{}, err
	}

	return UploadObjectResult{}, nil
}

type UploadObjectHandler decorator.Handler[UploadObjectArgs, UploadObjectResult]

func NewUploadObjectHandler(
	endpoint string,
	driveService remote.DriveService,
	objectRepo repository.ObjectRepository,
	eventRepo repository.EventRepository,
	capacityRepo repository.CapacityRepository,
	logger *logrus.Entry,
) UploadObjectHandler {
	return decorator.WithLogging[UploadObjectArgs, UploadObjectResult](
		uploadObjectHandler{
			endpoint:     endpoint,
			driveService: driveService,
			objectRepo:   objectRepo,
			eventRepo:    eventRepo,
			capacityRepo: capacityRepo,
		},
		logger,
	)
}
