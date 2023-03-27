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
	ShareFolderArgs struct {
		SessionId string
		FolderId  uuid.UUID
	}

	ShareFolderResult struct {
	}
)

type shareFolderHandler struct {
	userService remote.UserService
	driveRepo   repository.DriveRepository
	folderRepo  repository.FolderRepository
}

func (handler shareFolderHandler) Handle(ctx context.Context, args ShareFolderArgs) (ShareFolderResult, error) {
	user, err := handler.userService.GetUser(ctx, args.SessionId)
	if err != nil {
		return ShareFolderResult{}, err
	}

	drive, err := handler.driveRepo.GetDriveByOwner(ctx, user.Account())
	if err != nil {
		return ShareFolderResult{}, err
	}

	folder, err := handler.folderRepo.GetFolder(ctx, args.FolderId)
	if err != nil {
		return ShareFolderResult{}, err
	}

	err = folder.CanWrite(drive.Id())
	if err != nil {
		return ShareFolderResult{}, err
	}

	err = folder.Share()
	if err != nil {
		return ShareFolderResult{}, err
	}

	err = handler.folderRepo.UpdateFolder(ctx, folder, repository.UpdateFolderOptions{UpdateChildrenState: true})
	if err != nil {
		return ShareFolderResult{}, err
	}

	return ShareFolderResult{}, nil
}

type ShareFolderHandler decorator.Handler[ShareFolderArgs, ShareFolderResult]

func NewShareFolderHandler(
	userService remote.UserService,
	driveRepo repository.DriveRepository,
	folderRepo repository.FolderRepository,
	logger *logrus.Entry,
) ShareFolderHandler {
	return decorator.WithLogging[ShareFolderArgs, ShareFolderResult](
		shareFolderHandler{
			userService: userService,
			driveRepo:   driveRepo,
			folderRepo:  folderRepo,
		},
		logger,
	)
}
