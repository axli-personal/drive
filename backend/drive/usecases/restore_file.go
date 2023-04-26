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
	RestoreFileArgs struct {
		SessionId string
		FileId    uuid.UUID
	}

	RestoreFileResult struct {
	}
)

type restoreFileHandler struct {
	repo        repository.Repository
	userService remote.UserService
}

func (h restoreFileHandler) Handle(ctx context.Context, args RestoreFileArgs) (RestoreFileResult, error) {
	user, err := h.userService.GetUser(ctx, args.SessionId)
	if err != nil {
		return RestoreFileResult{}, err
	}

	drive, err := h.repo.GetDriveRepo().GetDriveByOwner(ctx, user.Account())
	if err != nil {
		return RestoreFileResult{}, err
	}

	file, err := h.repo.GetFileRepo().GetFile(ctx, args.FileId)
	if err != nil {
		return RestoreFileResult{}, err
	}

	err = file.CanWrite(drive.Id())
	if err != nil {
		return RestoreFileResult{}, err
	}

	err = file.ManuallyRestore()
	if err != nil {
		return RestoreFileResult{}, err
	}

	err = h.repo.GetFileRepo().UpdateFile(ctx, file, repository.UpdateFileOptions{})
	if err != nil {
		return RestoreFileResult{}, err
	}

	return RestoreFileResult{}, nil
}

type RestoreFileHandler decorator.Handler[RestoreFileArgs, RestoreFileResult]

func NewRestoreFileHandler(
	repo repository.Repository,
	userService remote.UserService,
	logger *logrus.Entry,
) RestoreFileHandler {
	return decorator.WithLogging[RestoreFileArgs, RestoreFileResult](
		restoreFileHandler{
			repo:        repo,
			userService: userService,
		},
		logger,
	)
}
