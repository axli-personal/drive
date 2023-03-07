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
	CanDownloadArgs struct {
		SessionId string
		FileId    uuid.UUID
	}

	CanDownloadResult struct {
		CanDownload bool
		FileName    string
	}
)

type canDownloadHandler struct {
	userService remote.UserService
	driveRepo   repository.DriveRepository
	fileRepo    repository.FileRepository
}

func (handler canDownloadHandler) Handle(ctx context.Context, args CanDownloadArgs) (CanDownloadResult, error) {
	user, err := handler.userService.GetUser(ctx, args.SessionId)
	if err != nil {
		return CanDownloadResult{}, err
	}

	drive, err := handler.driveRepo.GetDriveByOwner(ctx, user.Account())
	if err != nil {
		return CanDownloadResult{}, err
	}

	file, err := handler.fileRepo.GetFile(ctx, args.FileId)
	if err != nil {
		return CanDownloadResult{}, err
	}

	err = file.CanRead(drive.Id())
	if err != nil {
		return CanDownloadResult{}, err
	}

	return CanDownloadResult{
		CanDownload: true,
		FileName:    file.Name(),
	}, nil
}

type CanDownloadHandler decorator.Handler[CanDownloadArgs, CanDownloadResult]

func NewCanDownloadHandler(
	userService remote.UserService,
	driveRepo repository.DriveRepository,
	fileRepo repository.FileRepository,
	logger *logrus.Entry,
) CanDownloadHandler {
	return decorator.WithLogging[CanDownloadArgs, CanDownloadResult](
		canDownloadHandler{
			userService: userService,
			driveRepo:   driveRepo,
			fileRepo:    fileRepo,
		},
		logger,
	)
}
