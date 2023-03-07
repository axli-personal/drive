package usecases

import (
	"context"
	"github.com/axli-personal/drive/internal/common/decorator"
	"github.com/axli-personal/drive/internal/drive/domain"
	"github.com/axli-personal/drive/internal/drive/remote"
	"github.com/axli-personal/drive/internal/drive/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type (
	GetDriveArgs struct {
		SessionId string
	}

	GetDriveResult struct {
		Id       uuid.UUID
		Children Children
		Usage    domain.StorageUsage
		Plan     domain.StoragePlan
	}
)

type getDriveHandler struct {
	userService remote.UserService
	driveRepo   repository.DriveRepository
	folderRepo  repository.FolderRepository
	fileRepo    repository.FileRepository
}

func (handler getDriveHandler) Handle(ctx context.Context, args GetDriveArgs) (GetDriveResult, error) {
	user, err := handler.userService.GetUser(ctx, args.SessionId)
	if err != nil {
		return GetDriveResult{}, err
	}

	drive, err := handler.driveRepo.GetDriveByOwner(ctx, user.Account())
	if err != nil {
		return GetDriveResult{}, err
	}

	folders, err := handler.folderRepo.FindFolder(
		ctx,
		repository.FindFolderOptions{
			Parent: domain.CreateDriveParent(),
		},
	)

	files, err := handler.fileRepo.FindFile(
		ctx,
		repository.FindFileOptions{
			Parent: domain.CreateDriveParent(),
		},
	)

	return GetDriveResult{
		Id:       drive.Id(),
		Children: ToChildren(folders, files),
		Usage:    drive.Usage(),
		Plan:     drive.Plan(),
	}, nil
}

type GetDriveHandler decorator.Handler[GetDriveArgs, GetDriveResult]

func NewGetDriveHandler(
	userService remote.UserService,
	driveRepo repository.DriveRepository,
	folderRepo repository.FolderRepository,
	fileRepo repository.FileRepository,
	logger *logrus.Entry,
) GetDriveHandler {
	return decorator.WithLogging[GetDriveArgs, GetDriveResult](
		getDriveHandler{
			userService: userService,
			driveRepo:   driveRepo,
			folderRepo:  folderRepo,
			fileRepo:    fileRepo,
		},
		logger,
	)
}
