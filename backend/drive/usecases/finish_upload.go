package usecases

import (
	"context"
	"github.com/axli-personal/drive/backend/common/decorator"
	"github.com/axli-personal/drive/backend/drive/domain"
	"github.com/axli-personal/drive/backend/drive/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type FinishUploadArgs struct {
	FileId uuid.UUID
}

type FinishUploadResult struct {
}

type finishUploadHandler struct {
	fileRepo repository.FileRepository
}

func (handler finishUploadHandler) Handle(ctx context.Context, args FinishUploadArgs) (FinishUploadResult, error) {
	file, err := handler.fileRepo.GetFile(ctx, args.FileId)
	if err != nil {
		return FinishUploadResult{}, err
	}

	file.SetState(domain.StatePrivate)

	err = handler.fileRepo.UpdateFile(
		ctx,
		file,
		repository.UpdateFileOptions{
			MustInState: domain.StateLocked,
		},
	)
	if err != nil {
		return FinishUploadResult{}, err
	}

	return FinishUploadResult{}, err
}

type FinishUploadHandler decorator.Handler[FinishUploadArgs, FinishUploadResult]

func NewFileUploadedHandler(
	fileRepo repository.FileRepository,
	logger *logrus.Entry,
) FinishUploadHandler {
	return decorator.WithLogging[FinishUploadArgs, FinishUploadResult](
		finishUploadHandler{
			fileRepo: fileRepo,
		},
		logger,
	)
}
