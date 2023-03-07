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
	MoveFolderArgs struct {
		SessionId string
		FolderId  uuid.UUID
		Parent    domain.Parent
	}

	MoveFolderResult struct {
	}
)

type moveFolderHandler struct {
	userService remote.UserService
	driveRepo   repository.DriveRepository
	folderRepo  repository.FolderRepository
}

func (handler moveFolderHandler) Handle(
	ctx context.Context,
	args MoveFolderArgs,
) (MoveFolderResult, error) {
	user, err := handler.userService.GetUser(ctx, args.SessionId)
	if err != nil {
		return MoveFolderResult{}, err
	}

	drive, err := handler.driveRepo.GetDriveByOwner(ctx, user.Account())
	if err != nil {
		return MoveFolderResult{}, err
	}

	folder, err := handler.folderRepo.GetFolder(ctx, args.FolderId)
	if err != nil {
		return MoveFolderResult{}, err
	}

	err = folder.CanWrite(drive.Id())
	if err != nil {
		return MoveFolderResult{}, err
	}

	if args.Parent.IsFolder() {
		parentFolder, err := handler.folderRepo.GetFolder(ctx, args.Parent.FolderId())
		if err != nil {
			// Move to folder that does not exist.
			return MoveFolderResult{}, err
		}

		err = parentFolder.CanWrite(drive.Id())
		if err != nil {
			// Cannot write to parent folder.
			return MoveFolderResult{}, err
		}
	}

	err = folder.ChangeParent(args.Parent)
	if err != nil {
		return MoveFolderResult{}, err
	}

	// Only once transaction.
	err = handler.folderRepo.UpdateFolder(
		ctx,
		folder,
		repository.UpdateFolderOptions{
			MustInSameState: true,
		},
	)
	if err != nil {
		return MoveFolderResult{}, err
	}

	return MoveFolderResult{}, err
}

type MoveFolderHandler decorator.Handler[MoveFolderArgs, MoveFolderResult]

func NewMoveFolderHandler(
	userService remote.UserService,
	driveRepo repository.DriveRepository,
	folderRepo repository.FolderRepository,
	logger *logrus.Entry,
) MoveFolderHandler {
	return decorator.WithLogging[MoveFolderArgs, MoveFolderResult](
		moveFolderHandler{
			userService: userService,
			driveRepo:   driveRepo,
			folderRepo:  folderRepo,
		},
		logger,
	)
}
