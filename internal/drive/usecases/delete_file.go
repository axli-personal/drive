package usecases

import (
	"context"
	"github.com/axli-personal/drive/internal/common/decorator"
	"github.com/axli-personal/drive/internal/drive/remote"
	"github.com/axli-personal/drive/internal/drive/repository"
	"github.com/axli-personal/drive/internal/pkg/events"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type (
	DeleteFileArgs struct {
		SessionId string
		FileId    uuid.UUID
	}

	DeleteFileResult struct {
	}
)

type deleteFileHandler struct {
	userService remote.UserService
	driveRepo   repository.DriveRepository
	fileRepo    repository.FileRepository
	eventRepo   repository.EventRepository
}

func (handler deleteFileHandler) Handle(
	ctx context.Context,
	args DeleteFileArgs,
) (DeleteFileResult, error) {
	user, err := handler.userService.GetUser(ctx, args.SessionId)
	if err != nil {
		return DeleteFileResult{}, err
	}

	drive, err := handler.driveRepo.GetDriveByOwner(ctx, user.Account())
	if err != nil {
		return DeleteFileResult{}, err
	}

	file, err := handler.fileRepo.GetFile(ctx, args.FileId)
	if err != nil {
		return DeleteFileResult{}, err
	}

	err = file.CanDelete(drive.Id())
	if err != nil {
		return DeleteFileResult{}, err
	}

	err = handler.fileRepo.DeleteFile(ctx, file)
	if err != nil {
		return DeleteFileResult{}, err
	}

	// publish may fail.
	err = handler.eventRepo.PublishFileDeleted(
		ctx,
		events.FileDeleted{
			FileId: file.Id().String(),
		},
	)

	return DeleteFileResult{}, nil
}

type DeleteFileHandler decorator.Handler[DeleteFileArgs, DeleteFileResult]

func NewDeleteFileHandler(
	userService remote.UserService,
	driveRepo repository.DriveRepository,
	fileRepo repository.FileRepository,
	eventRepo repository.EventRepository,
	logger *logrus.Entry,
) DeleteFileHandler {
	return decorator.WithLogging[DeleteFileArgs, DeleteFileResult](
		deleteFileHandler{
			userService: userService,
			driveRepo:   driveRepo,
			fileRepo:    fileRepo,
			eventRepo:   eventRepo,
		},
		logger,
	)
}
