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
	RemoveFolderArgs struct {
		SessionId string
		FolderId  uuid.UUID
	}

	RemoveFolderResult struct {
	}
)

type removeFolderHandler struct {
	userService remote.UserService
	driveRepo   repository.DriveRepository
	folderRepo  repository.FolderRepository
}

func (handler removeFolderHandler) Handle(ctx context.Context, args RemoveFolderArgs) (RemoveFolderResult, error) {
	user, err := handler.userService.GetUser(ctx, args.SessionId)
	if err != nil {
		return RemoveFolderResult{}, err
	}

	drive, err := handler.driveRepo.GetDriveByOwner(ctx, user.Account())
	if err != nil {
		return RemoveFolderResult{}, err
	}

	folder, err := handler.folderRepo.GetFolder(ctx, args.FolderId)
	if err != nil {
		return RemoveFolderResult{}, err
	}

	err = folder.CanWrite(drive.Id())
	if err != nil {
		return RemoveFolderResult{}, err
	}

	err = folder.Trash()
	if err != nil {
		return RemoveFolderResult{}, err
	}

	err = handler.folderRepo.UpdateFolder(
		ctx,
		folder,
		repository.UpdateFolderOptions{
			UpdateChildrenState: true,
		},
	)
	if err != nil {
		return RemoveFolderResult{}, err
	}

	return RemoveFolderResult{}, nil
}

type RemoveFolderHandler decorator.Handler[RemoveFolderArgs, RemoveFolderResult]

func NewRemoveFolderHandler(
	userService remote.UserService,
	driveRepo repository.DriveRepository,
	folderRepo repository.FolderRepository,
	logger *logrus.Entry,
) RemoveFolderHandler {
	return decorator.WithLogging[RemoveFolderArgs, RemoveFolderResult](
		removeFolderHandler{
			userService: userService,
			driveRepo:   driveRepo,
			folderRepo:  folderRepo,
		},
		logger,
	)
}
