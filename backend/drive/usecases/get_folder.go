package usecases

import (
	"context"
	"github.com/axli-personal/drive/backend/common/decorator"
	"github.com/axli-personal/drive/backend/drive/domain"
	"github.com/axli-personal/drive/backend/drive/remote"
	"github.com/axli-personal/drive/backend/drive/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"time"
)

type (
	GetFolderArgs struct {
		FolderId  uuid.UUID
		SessionId string
	}

	GetFolderResult struct {
		FolderId   uuid.UUID
		Parent     domain.Parent
		Name       string
		Shared     bool
		LastChange time.Time
		Children   Children
	}
)

type getFolderHandler struct {
	userService remote.UserService
	driveRepo   repository.DriveRepository
	folderRepo  repository.FolderRepository
	fileRepo    repository.FileRepository
}

func (handler getFolderHandler) Handle(ctx context.Context, args GetFolderArgs) (GetFolderResult, error) {
	userDriveId := uuid.Nil

	if args.SessionId != "" {
		user, err := handler.userService.GetUser(ctx, args.SessionId)
		if err != nil {
			return GetFolderResult{}, err
		}

		drive, err := handler.driveRepo.GetDriveByOwner(ctx, user.Account())
		if err != nil {
			return GetFolderResult{}, err
		}

		userDriveId = drive.Id()
	}

	folder, err := handler.folderRepo.GetFolder(ctx, args.FolderId)
	if err != nil {
		return GetFolderResult{}, err
	}

	err = folder.CanRead(userDriveId)
	if err != nil {
		return GetFolderResult{}, err
	}

	result := GetFolderResult{
		FolderId:   folder.Id(),
		Name:       folder.Name(),
		Shared:     folder.State() == domain.StateShared,
		LastChange: folder.LastChange(),
	}

	self, err := domain.CreateFolderParent(args.FolderId)
	if err != nil {
		return result, err
	}

	folders, err := handler.folderRepo.FindFolder(
		ctx,
		repository.FindFolderOptions{
			Parent: self,
		},
	)
	if err != nil {
		return result, err
	}

	files, err := handler.fileRepo.FindFile(
		ctx,
		repository.FindFileOptions{
			Parent: self,
		},
	)
	if err != nil {
		return result, err
	}

	result.Children = ToChildren(folders, files)

	if folder.CanReadParent(userDriveId) == nil {
		result.Parent = folder.Parent()
	}

	return result, nil
}

type GetFolderHandler decorator.Handler[GetFolderArgs, GetFolderResult]

func NewGetFolderHandler(
	userService remote.UserService,
	driveRepo repository.DriveRepository,
	folderRepo repository.FolderRepository,
	fileRepo repository.FileRepository,
	logger *logrus.Entry,
) GetFolderHandler {
	return decorator.WithLogging[GetFolderArgs, GetFolderResult](
		getFolderHandler{
			userService: userService,
			driveRepo:   driveRepo,
			folderRepo:  folderRepo,
			fileRepo:    fileRepo,
		},
		logger,
	)
}
