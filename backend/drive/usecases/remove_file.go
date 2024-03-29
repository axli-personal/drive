package usecases

import (
	"context"
	"github.com/axli-personal/drive/backend/common/decorator"
	"github.com/axli-personal/drive/backend/drive/remote"
	"github.com/axli-personal/drive/backend/drive/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type (
	RemoveFileArgs struct {
		SessionId string
		FileId    uuid.UUID
	}

	RemoveFileResult struct {
	}
)

type removeFileHandler struct {
	repo        repository.Repository
	userService remote.UserService
}

func (h removeFileHandler) Handle(ctx context.Context, args RemoveFileArgs) (RemoveFileResult, error) {
	user, err := h.userService.GetUser(ctx, args.SessionId)
	if err != nil {
		return RemoveFileResult{}, err
	}

	drive, err := h.repo.GetDriveRepo().GetDriveByOwner(ctx, user.Account())
	if err != nil {
		return RemoveFileResult{}, err
	}

	file, err := h.repo.GetFileRepo().GetFile(ctx, args.FileId)
	if err != nil {
		return RemoveFileResult{}, err
	}

	err = file.CanWrite(drive.Id())
	if err != nil {
		return RemoveFileResult{}, err
	}

	err = file.ManuallyTrash()
	if err != nil {
		return RemoveFileResult{}, err
	}

	err = h.repo.GetFileRepo().UpdateFile(ctx, file, repository.UpdateFileOptions{})
	if err != nil {
		return RemoveFileResult{}, err
	}

	return RemoveFileResult{}, nil
}

type RemoveFileHandler decorator.Handler[RemoveFileArgs, RemoveFileResult]

func NewRemoveFileHandler(
	repo repository.Repository,
	userService remote.UserService,
	logger *logrus.Entry,
) RemoveFileHandler {
	return decorator.WithLogging[RemoveFileArgs, RemoveFileResult](
		removeFileHandler{
			repo:        repo,
			userService: userService,
		},
		logger,
	)
}
