package usecases

import (
	"context"
	"github.com/axli-personal/drive/backend/common/decorator"
	"github.com/axli-personal/drive/backend/pkg/types"
	"github.com/axli-personal/drive/backend/storage/remote"
	"github.com/axli-personal/drive/backend/storage/repository"
	"github.com/sirupsen/logrus"
	"io"
)

type (
	DownloadArgs struct {
		SessionId string
		FileId    string
	}

	DownloadResult struct {
		Data     io.Reader
		FileName string
	}
)

type downloadHandler struct {
	driveService remote.DriveService
	objectRepo   repository.ObjectRepository
}

func (handler downloadHandler) Handle(ctx context.Context, args DownloadArgs) (DownloadResult, error) {
	response, err := handler.driveService.StartDownload(
		types.StartDownloadRequest{
			SessionId: args.SessionId,
			FileId:    args.FileId,
		},
	)
	if err != nil {
		return DownloadResult{}, err
	}

	object, err := handler.objectRepo.GetObject(ctx, response.FileHash)
	if err != nil {
		return DownloadResult{}, err
	}

	return DownloadResult{
		FileName: response.FileName,
		Data:     object,
	}, nil
}

type DownloadHandler decorator.Handler[DownloadArgs, DownloadResult]

func NewDownloadHandler(
	driveService remote.DriveService,
	objectRepo repository.ObjectRepository,
	logger *logrus.Entry,
) DownloadHandler {
	return decorator.WithLogging[DownloadArgs, DownloadResult](
		downloadHandler{
			driveService: driveService,
			objectRepo:   objectRepo,
		},
		logger,
	)
}
