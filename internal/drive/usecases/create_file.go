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
	CreateFileArgs struct {
		SessionId  string
		FileParent domain.Parent
		FileName   string
	}

	CreateFileResult struct {
		FileId          uuid.UUID
		FileName        string
		LastChange      time.Time
		StorageEndpoint string
	}
)

type createFileHandler struct {
	userService    remote.UserService
	storageCluster remote.StorageCluster
	driveRepo      repository.DriveRepository
	fileRepo       repository.FileRepository
}

func (handler createFileHandler) Handle(ctx context.Context, args CreateFileArgs) (CreateFileResult, error) {
	user, err := handler.userService.GetUser(ctx, args.SessionId)
	if err != nil {
		return CreateFileResult{}, err
	}

	drive, err := handler.driveRepo.GetDriveByOwner(ctx, user.Account())
	if err != nil {
		return CreateFileResult{}, err
	}

	file, err := domain.NewFile(drive.Id(), args.FileParent, args.FileName)
	if err != nil {
		return CreateFileResult{}, err
	}

	err = handler.fileRepo.SaveFile(ctx, file)
	if err != nil {
		return CreateFileResult{}, err
	}

	endpoint, err := handler.storageCluster.ChooseStorageEndPoint(ctx)
	if err != nil {
		return CreateFileResult{}, err
	}

	return CreateFileResult{
		FileId:          file.Id(),
		FileName:        file.Name(),
		LastChange:      file.LastChange(),
		StorageEndpoint: endpoint,
	}, nil
}

type CreateFileHandler decorator.Handler[CreateFileArgs, CreateFileResult]

func NewCreateFileHandler(
	userService remote.UserService,
	storageCluster remote.StorageCluster,
	driveRepo repository.DriveRepository,
	fileRepo repository.FileRepository,
	logger *logrus.Entry,
) CreateFileHandler {
	return decorator.WithLogging[CreateFileArgs, CreateFileResult](
		createFileHandler{
			userService:    userService,
			storageCluster: storageCluster,
			driveRepo:      driveRepo,
			fileRepo:       fileRepo,
		},
		logger,
	)
}
