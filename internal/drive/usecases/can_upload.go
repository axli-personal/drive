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
	CanUploadArgs struct {
		SessionId string
		FileId    uuid.UUID
	}

	CanUploadResult struct {
		CanUpload bool
	}
)

type canUploadHandler struct {
	userService remote.UserService
	driveRepo   repository.DriveRepository
	fileRepo    repository.FileRepository
}

func (handler canUploadHandler) Handle(ctx context.Context, args CanUploadArgs) (CanUploadResult, error) {
	user, err := handler.userService.GetUser(ctx, args.SessionId)
	if err != nil {
		return CanUploadResult{}, err
	}

	drive, err := handler.driveRepo.GetDriveByOwner(ctx, user.Account())
	if err != nil {
		return CanUploadResult{}, err
	}

	file, err := handler.fileRepo.GetFile(ctx, args.FileId)
	if err != nil {
		return CanUploadResult{}, err
	}

	err = file.CanWrite(drive.Id())
	if err != nil {
		return CanUploadResult{}, err
	}

	return CanUploadResult{CanUpload: true}, nil
}

type CanUploadHandler decorator.Handler[CanUploadArgs, CanUploadResult]

func NewCanUploadHandler(
	userService remote.UserService,
	driveRepo repository.DriveRepository,
	fileRepo repository.FileRepository,
	logger *logrus.Entry,
) CanUploadHandler {
	return decorator.WithLogging[CanUploadArgs, CanUploadResult](
		canUploadHandler{
			userService: userService,
			driveRepo:   driveRepo,
			fileRepo:    fileRepo,
		},
		logger,
	)
}
