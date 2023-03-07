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
	CreateFolderArgs struct {
		SessionId    string
		FolderParent domain.Parent
		FolderName   string
	}

	CreateFolderResult struct {
		FolderId   uuid.UUID
		FolderName string
		LastChange time.Time
	}
)

type createFolderHandler struct {
	userService remote.UserService
	driveRepo   repository.DriveRepository
	folderRepo  repository.FolderRepository
}

func (handler createFolderHandler) Handle(ctx context.Context, args CreateFolderArgs) (CreateFolderResult, error) {
	user, err := handler.userService.GetUser(ctx, args.SessionId)
	if err != nil {
		return CreateFolderResult{}, err
	}

	drive, err := handler.driveRepo.GetDriveByOwner(ctx, user.Account())
	if err != nil {
		return CreateFolderResult{}, err
	}

	folder, err := domain.NewFolder(drive.Id(), args.FolderParent, args.FolderName)
	if err != nil {
		return CreateFolderResult{}, err
	}

	err = handler.folderRepo.SaveFolder(ctx, folder)
	if err != nil {
		return CreateFolderResult{}, err
	}

	return CreateFolderResult{
		FolderId:   folder.Id(),
		FolderName: folder.Name(),
		LastChange: folder.LastChange(),
	}, nil
}

type CreateFolderHandler decorator.Handler[CreateFolderArgs, CreateFolderResult]

func NewCreateFolderHandler(
	userService remote.UserService,
	driveRepo repository.DriveRepository,
	folderRepo repository.FolderRepository,
	logger *logrus.Entry,
) CreateFolderHandler {
	return decorator.WithLogging[CreateFolderArgs, CreateFolderResult](
		createFolderHandler{
			userService: userService,
			driveRepo:   driveRepo,
			folderRepo:  folderRepo,
		},
		logger,
	)
}
