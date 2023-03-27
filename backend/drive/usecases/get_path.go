package usecases

import (
	"context"
	"github.com/axli-personal/drive/backend/common/decorator"
	"github.com/axli-personal/drive/backend/drive/domain"
	"github.com/axli-personal/drive/backend/drive/remote"
	"github.com/axli-personal/drive/backend/drive/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type GetPathArgs struct {
	SessionId string
	Parent    domain.Parent
}

type GetPathResult struct {
	Folders []FolderLink
}

type getPathHandler struct {
	userService remote.UserService
	driveRepo   repository.DriveRepository
	folderRepo  repository.FolderRepository
}

func (handler getPathHandler) Handle(ctx context.Context, args GetPathArgs) (GetPathResult, error) {
	userDriveId := uuid.Nil

	if args.SessionId != "" {
		user, err := handler.userService.GetUser(ctx, args.SessionId)
		if err != nil {
			return GetPathResult{}, err
		}

		drive, err := handler.driveRepo.GetDriveByOwner(ctx, user.Account())
		if err != nil {
			return GetPathResult{}, err
		}

		userDriveId = drive.Id()
	}

	result := GetPathResult{}

	currentParent := args.Parent
	for currentParent.IsFolder() {
		FolderId := currentParent.FolderId()

		folder, err := handler.folderRepo.GetFolder(ctx, FolderId)
		if err != nil {
			break
		}

		err = folder.CanRead(userDriveId)
		if err != nil {
			break
		}

		result.Folders = append(result.Folders, FolderLink{
			Id:   folder.Id(),
			Name: folder.Name(),
		})

		currentParent = folder.Parent()
	}

	return result, nil
}

type GetPathHandler decorator.Handler[GetPathArgs, GetPathResult]

func NewGetPathHandler(
	userService remote.UserService,
	driveRepo repository.DriveRepository,
	folderRepo repository.FolderRepository,
	logger *logrus.Entry,
) GetPathHandler {
	return decorator.WithLogging[GetPathArgs, GetPathResult](
		getPathHandler{
			userService: userService,
			driveRepo:   driveRepo,
			folderRepo:  folderRepo,
		},
		logger,
	)
}
