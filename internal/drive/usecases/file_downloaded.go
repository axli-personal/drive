package usecases

import (
	"context"
	"github.com/axli-personal/drive/internal/common/decorator"
	"github.com/axli-personal/drive/internal/drive/repository"
	"github.com/axli-personal/drive/internal/pkg/events"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type FileDownloadedResult struct {
	CanAcknowledge bool
}

type fileDownloadedHandler struct {
	fileRepo repository.FileRepository
}

func (handler fileDownloadedHandler) Handle(
	ctx context.Context,
	event events.FileDownloaded,
) (FileDownloadedResult, error) {
	fileId, err := uuid.Parse(event.FileId)
	if err != nil {
		return FileDownloadedResult{CanAcknowledge: true}, err
	}

	file, err := handler.fileRepo.GetFile(ctx, fileId)
	if err != nil {
		if err == repository.ErrNotFound {
			return FileDownloadedResult{CanAcknowledge: true}, err
		} else {
			// Some unknown error, need retry later.
			return FileDownloadedResult{CanAcknowledge: false}, err
		}
	}

	err = file.IncreaseDownloadTimes()
	if err != nil {
		return FileDownloadedResult{CanAcknowledge: true}, err
	}

	// Only once transaction.
	err = handler.fileRepo.UpdateFile(
		ctx,
		file,
		repository.UpdateFileOptions{
			MustInState: file.State(),
		},
	)
	if err != nil {
		return FileDownloadedResult{CanAcknowledge: false}, err
	}

	return FileDownloadedResult{CanAcknowledge: true}, nil
}

type FileDownloadedHandler decorator.Handler[events.FileDownloaded, FileDownloadedResult]

func NewFileDownloadedHandler(
	fileRepo repository.FileRepository,
	logger *logrus.Entry,
) FileDownloadedHandler {
	return decorator.WithLogging[events.FileDownloaded, FileDownloadedResult](
		fileDownloadedHandler{
			fileRepo: fileRepo,
		},
		logger,
	)
}
