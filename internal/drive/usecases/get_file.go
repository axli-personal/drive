package usecases

import (
	"context"
	"github.com/axli-personal/drive/internal/common/decorator"
	"github.com/axli-personal/drive/internal/drive/domain"
	"github.com/axli-personal/drive/internal/drive/remote"
	"github.com/axli-personal/drive/internal/drive/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"time"
)

type (
	GetFileArgs struct {
		FileId    uuid.UUID
		SessionId string
	}

	GetFileResult struct {
		FileId         uuid.UUID
		Parent         domain.Parent
		Name           string
		LastChange     time.Time
		Endpoint       string
		Bytes          int
		downloadCounts int
	}
)

type getFileHandler struct {
	userService remote.UserService
	driveRepo   repository.DriveRepository
	fileRepo    repository.FileRepository
	logger      *logrus.Entry
}

func (handler getFileHandler) Handle(ctx context.Context, args GetFileArgs) (GetFileResult, error) {
	userDriveId := uuid.Nil

	if args.SessionId != "" {
		user, err := handler.userService.GetUser(ctx, args.SessionId)
		if err != nil {
			return GetFileResult{}, err
		}

		drive, err := handler.driveRepo.GetDriveByOwner(ctx, user.Account())
		if err != nil {
			return GetFileResult{}, err
		}

		userDriveId = drive.Id()
	}

	file, err := handler.fileRepo.GetFile(ctx, args.FileId)
	if err != nil {
		return GetFileResult{}, err
	}

	err = file.CanRead(userDriveId)
	if err != nil {
		return GetFileResult{}, err
	}

	result := GetFileResult{
		FileId:         file.Id(),
		Name:           file.Name(),
		LastChange:     file.LastChange(),
		Endpoint:       file.Endpoint(),
		Bytes:          file.Bytes(),
		downloadCounts: file.DownloadCounts(),
	}

	if file.CanReadParent(userDriveId) == nil {
		result.Parent = file.Parent()
	}

	return result, nil
}

type GetFileHandler decorator.Handler[GetFileArgs, GetFileResult]

func NewGetFileHandler(
	userService remote.UserService,
	driveRepo repository.DriveRepository,
	fileRepo repository.FileRepository,
	logger *logrus.Entry,
) GetFileHandler {
	return decorator.WithLogging[GetFileArgs, GetFileResult](
		getFileHandler{
			userService: userService,
			driveRepo:   driveRepo,
			fileRepo:    fileRepo,
			logger:      logger,
		},
		logger,
	)
}
