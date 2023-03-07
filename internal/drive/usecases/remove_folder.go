package usecases

import (
	"context"
	"github.com/axli-personal/drive/internal/common/decorator"
	"github.com/axli-personal/drive/internal/drive/remote"
	"github.com/axli-personal/drive/internal/drive/repository"
	"github.com/axli-personal/drive/internal/pkg/events"
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
	eventRepo   repository.EventRepository
}

func (handler removeFolderHandler) Handle(
	ctx context.Context,
	args RemoveFolderArgs,
) (RemoveFolderResult, error) {
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

	oldState := folder.State()

	err = folder.Trash()
	if err != nil {
		return RemoveFolderResult{}, err
	}

	// Only once transaction.
	err = handler.folderRepo.UpdateFolder(
		ctx,
		folder,
		repository.UpdateFolderOptions{
			MustInState: oldState,
		},
	)
	if err != nil {
		return RemoveFolderResult{}, err
	}

	// TODO: publish may fail.
	err = handler.eventRepo.PublishFolderRemoved(
		ctx,
		events.FolderRemoved{
			FolderId: folder.Id().String(),
		},
	)

	return RemoveFolderResult{}, nil
}

type RemoveFolderHandler decorator.Handler[RemoveFolderArgs, RemoveFolderResult]

func NewRemoveFolderHandler(
	userService remote.UserService,
	driveRepo repository.DriveRepository,
	folderRepo repository.FolderRepository,
	eventRepo repository.EventRepository,
	logger *logrus.Entry,
) RemoveFolderHandler {
	return decorator.WithLogging[RemoveFolderArgs, RemoveFolderResult](
		removeFolderHandler{
			userService: userService,
			driveRepo:   driveRepo,
			folderRepo:  folderRepo,
			eventRepo:   eventRepo,
		},
		logger,
	)
}
