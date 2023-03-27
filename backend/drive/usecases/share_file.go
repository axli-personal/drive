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
	ShareFileArgs struct {
		SessionId string
		FileId    uuid.UUID
	}

	ShareFileResult struct {
	}
)

type shareFileHandler struct {
	userService remote.UserService
	driveRepo   repository.DriveRepository
	fileRepo    repository.FileRepository
}

func (handler shareFileHandler) Handle(ctx context.Context, args ShareFileArgs) (ShareFileResult, error) {
	user, err := handler.userService.GetUser(ctx, args.SessionId)
	if err != nil {
		return ShareFileResult{}, err
	}

	drive, err := handler.driveRepo.GetDriveByOwner(ctx, user.Account())
	if err != nil {
		return ShareFileResult{}, err
	}

	file, err := handler.fileRepo.GetFile(ctx, args.FileId)
	if err != nil {
		return ShareFileResult{}, err
	}

	err = file.CanWrite(drive.Id())
	if err != nil {
		return ShareFileResult{}, err
	}

	err = file.Share()
	if err != nil {
		return ShareFileResult{}, err
	}

	err = handler.fileRepo.UpdateFile(ctx, file, repository.UpdateFileOptions{})
	if err != nil {
		return ShareFileResult{}, err
	}

	return ShareFileResult{}, nil
}

type ShareFileHandler decorator.Handler[ShareFileArgs, ShareFileResult]

func NewShareFileHandler(
	userService remote.UserService,
	driveRepo repository.DriveRepository,
	fileRepo repository.FileRepository,
	logger *logrus.Entry,
) ShareFileHandler {
	return decorator.WithLogging[ShareFileArgs, ShareFileResult](
		shareFileHandler{
			userService: userService,
			driveRepo:   driveRepo,
			fileRepo:    fileRepo,
		},
		logger,
	)
}
