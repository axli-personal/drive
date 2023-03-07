package usecases

import (
	"context"
	"github.com/axli-personal/drive/internal/common/decorator"
	"github.com/axli-personal/drive/internal/drive/remote"
	"github.com/axli-personal/drive/internal/drive/repository"
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
	userService remote.UserService
	driveRepo   repository.DriveRepository
	fileRepo    repository.FileRepository
}

func (handler removeFileHandler) Handle(
	ctx context.Context,
	args RemoveFileArgs,
) (RemoveFileResult, error) {
	user, err := handler.userService.GetUser(ctx, args.SessionId)
	if err != nil {
		return RemoveFileResult{}, err
	}

	drive, err := handler.driveRepo.GetDriveByOwner(ctx, user.Account())
	if err != nil {
		return RemoveFileResult{}, err
	}

	file, err := handler.fileRepo.GetFile(ctx, args.FileId)
	if err != nil {
		return RemoveFileResult{}, err
	}

	err = file.CanWrite(drive.Id())
	if err != nil {
		return RemoveFileResult{}, err
	}

	oldState := file.State()

	err = file.Trash()
	if err != nil {
		return RemoveFileResult{}, err
	}

	// Only once transaction.
	err = handler.fileRepo.UpdateFile(
		ctx,
		file,
		repository.UpdateFileOptions{
			MustInState: oldState,
		},
	)
	if err != nil {
		return RemoveFileResult{}, err
	}

	return RemoveFileResult{}, nil
}

type RemoveFileHandler decorator.Handler[RemoveFileArgs, RemoveFileResult]

func NewRemoveFileHandler(
	userService remote.UserService,
	driveRepo repository.DriveRepository,
	fileRepo repository.FileRepository,
	logger *logrus.Entry,
) RemoveFileHandler {
	return decorator.WithLogging[RemoveFileArgs, RemoveFileResult](
		removeFileHandler{
			userService: userService,
			driveRepo:   driveRepo,
			fileRepo:    fileRepo,
		},
		logger,
	)
}
