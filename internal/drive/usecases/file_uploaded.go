package usecases

import (
	"context"
	"github.com/axli-personal/drive/internal/common/decorator"
	"github.com/axli-personal/drive/internal/drive/repository"
	"github.com/axli-personal/drive/internal/pkg/events"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type FileUploadedResult struct {
	CanAcknowledge bool
}

type fileUploadedHandler struct {
	fileRepo repository.FileRepository
}

func (handler fileUploadedHandler) Handle(
	ctx context.Context,
	event events.FileUploaded,
) (FileUploadedResult, error) {
	fileId, err := uuid.Parse(event.FileId)
	if err != nil {
		return FileUploadedResult{CanAcknowledge: true}, err
	}

	file, err := handler.fileRepo.GetFile(ctx, fileId)
	if err != nil {
		if err == repository.ErrNotFound {
			return FileUploadedResult{CanAcknowledge: true}, err
		} else {
			// Some unknown error, need retry later.
			return FileUploadedResult{CanAcknowledge: false}, err
		}
	}

	oldState := file.State()

	// File not in valid state, maybe bind twice.
	err = file.BindStorage(event.Endpoint, event.TotalBytes)
	if err != nil {
		return FileUploadedResult{CanAcknowledge: true}, err
	}

	// Only once transaction.
	err = handler.fileRepo.UpdateFile(
		ctx,
		file,
		repository.UpdateFileOptions{
			MustInState:          oldState,
			IncreaseStorageUsage: true,
		},
	)
	if err != nil {
		return FileUploadedResult{CanAcknowledge: false}, err
	}

	return FileUploadedResult{CanAcknowledge: true}, err
}

type FileUploadedHandler decorator.Handler[events.FileUploaded, FileUploadedResult]

func NewFileUploadedHandler(
	fileRepo repository.FileRepository,
	logger *logrus.Entry,
) FileUploadedHandler {
	return decorator.WithLogging[events.FileUploaded, FileUploadedResult](
		fileUploadedHandler{
			fileRepo: fileRepo,
		},
		logger,
	)
}
