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
	StartDownloadArgs struct {
		SessionId string
		FileId    uuid.UUID
	}

	StartDownloadResult struct {
		FileName string
		FileHash string
	}
)

type startDownloadHandler struct {
	userService remote.UserService
	driveRepo   repository.DriveRepository
	fileRepo    repository.FileRepository
}

func (handler startDownloadHandler) Handle(ctx context.Context, args StartDownloadArgs) (StartDownloadResult, error) {
	var driveId uuid.UUID

	if args.SessionId != "" {
		user, err := handler.userService.GetUser(ctx, args.SessionId)
		if err != nil {
			return StartDownloadResult{}, err
		}

		drive, err := handler.driveRepo.GetDriveByOwner(ctx, user.Account())
		if err != nil {
			return StartDownloadResult{}, err
		}

		driveId = drive.Id()
	}

	file, err := handler.fileRepo.GetFile(ctx, args.FileId)
	if err != nil {
		return StartDownloadResult{}, err
	}

	err = file.CanRead(driveId)
	if err != nil {
		return StartDownloadResult{}, err
	}

	file.IncreaseDownloadCounts()

	err = handler.fileRepo.UpdateFile(
		ctx,
		file,
		repository.UpdateFileOptions{},
	)
	if err != nil {
		return StartDownloadResult{}, err
	}

	return StartDownloadResult{
		FileName: file.Name(),
		FileHash: file.Hash(),
	}, nil
}

type StartDownloadHandler decorator.Handler[StartDownloadArgs, StartDownloadResult]

func NewStartDownloadHandler(
	userService remote.UserService,
	driveRepo repository.DriveRepository,
	fileRepo repository.FileRepository,
	logger *logrus.Entry,
) StartDownloadHandler {
	return decorator.WithLogging[StartDownloadArgs, StartDownloadResult](
		startDownloadHandler{
			userService: userService,
			driveRepo:   driveRepo,
			fileRepo:    fileRepo,
		},
		logger,
	)
}
