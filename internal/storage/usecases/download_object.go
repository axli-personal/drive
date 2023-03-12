package usecases

import (
	"context"
	"errors"
	"github.com/axli-personal/drive/internal/common/decorator"
	"github.com/axli-personal/drive/internal/pkg/events"
	"github.com/axli-personal/drive/internal/pkg/types"
	"github.com/axli-personal/drive/internal/storage/remote"
	"github.com/axli-personal/drive/internal/storage/repository"
	"github.com/sirupsen/logrus"
	"io"
)

var (
	ErrCannotDownload = errors.New("cannot download")
)

type (
	DownloadObjectArgs struct {
		SessionId string
		FileId    string
	}

	DownloadObjectResult struct {
		Data     io.Reader
		FileName string
	}
)

type downloadObjectHandler struct {
	driveService remote.DriveService
	objectRepo   repository.ObjectRepository
	eventRepo    repository.EventRepository
	capacityRepo repository.CapacityRepository
}

func (handler downloadObjectHandler) Handle(ctx context.Context, args DownloadObjectArgs) (DownloadObjectResult, error) {
	err := handler.capacityRepo.DecreaseRequestCapacity(ctx)
	if err != nil {
		return DownloadObjectResult{}, err
	}

	response, err := handler.driveService.CanDownload(
		types.CanDownloadRequest{
			SessionId: args.SessionId,
			FileId:    args.FileId,
		},
	)
	if err != nil {
		return DownloadObjectResult{}, err
	}
	if !response.CanDownload {
		return DownloadObjectResult{}, ErrCannotDownload
	}

	object, err := handler.objectRepo.GetObject(ctx, args.FileId)
	if err != nil {
		return DownloadObjectResult{}, err
	}

	err = handler.eventRepo.PublishFileDownloaded(
		ctx,
		events.FileDownloaded{
			FileId: object.FileId(),
		},
	)
	if err != nil {
		// TODO: log fail to publish events.
	}

	return DownloadObjectResult{
		Data:     object,
		FileName: response.FileName,
	}, nil
}

type DownloadObjectHandler decorator.Handler[DownloadObjectArgs, DownloadObjectResult]

func NewDownloadObjectHandler(
	driveService remote.DriveService,
	objectRepo repository.ObjectRepository,
	eventRepo repository.EventRepository,
	capacityRepo repository.CapacityRepository,
	logger *logrus.Entry,
) DownloadObjectHandler {
	return decorator.WithLogging[DownloadObjectArgs, DownloadObjectResult](
		downloadObjectHandler{
			driveService: driveService,
			objectRepo:   objectRepo,
			eventRepo:    eventRepo,
			capacityRepo: capacityRepo,
		},
		logger,
	)
}
