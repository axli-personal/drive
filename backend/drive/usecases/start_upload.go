package usecases

import (
	"context"
	"github.com/axli-personal/drive/backend/common/decorator"
	"github.com/axli-personal/drive/backend/drive/domain"
	"github.com/axli-personal/drive/backend/drive/remote"
	"github.com/axli-personal/drive/backend/drive/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type (
	StartUploadArgs struct {
		SessionId  string
		FileParent domain.Parent
		FileName   string
		FileHash   string
		FileSize   int
	}

	StartUploadResult struct {
		FileId uuid.UUID
	}
)

type startUploadHandler struct {
	userService remote.UserService
	driveRepo   repository.DriveRepository
	fileRepo    repository.FileRepository
}

func (handler startUploadHandler) Handle(ctx context.Context, args StartUploadArgs) (StartUploadResult, error) {
	user, err := handler.userService.GetUser(ctx, args.SessionId)
	if err != nil {
		return StartUploadResult{}, err
	}

	drive, err := handler.driveRepo.GetDriveByOwner(ctx, user.Account())
	if err != nil {
		return StartUploadResult{}, err
	}

	err = drive.IncreaseUsage(args.FileSize)
	if err != nil {
		return StartUploadResult{}, err
	}

	err = handler.driveRepo.UpdateDrive(ctx, drive)
	if err != nil {
		return StartUploadResult{}, err
	}

	file, err := domain.NewFile(drive.Id(), args.FileParent, args.FileName, args.FileSize, args.FileHash)
	if err != nil {
		return StartUploadResult{}, err
	}

	err = handler.fileRepo.SaveFile(ctx, file)
	if err != nil {
		return StartUploadResult{}, err
	}

	return StartUploadResult{FileId: file.Id()}, nil
}

type StartUploadHandler decorator.Handler[StartUploadArgs, StartUploadResult]

func NewStartUploadHandler(
	userService remote.UserService,
	driveRepo repository.DriveRepository,
	fileRepo repository.FileRepository,
	logger *logrus.Entry,
) StartUploadHandler {
	return decorator.WithLogging[StartUploadArgs, StartUploadResult](
		startUploadHandler{
			userService: userService,
			driveRepo:   driveRepo,
			fileRepo:    fileRepo,
		},
		logger,
	)
}
