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
	MoveFileArgs struct {
		SessionId string
		FileId    uuid.UUID
		Parent    domain.Parent
	}

	MoveFileResult struct {
	}
)

type moveFileHandler struct {
	userService remote.UserService
	driveRepo   repository.DriveRepository
	folderRepo  repository.FolderRepository
	fileRepo    repository.FileRepository
}

func (handler moveFileHandler) Handle(ctx context.Context, args MoveFileArgs) (MoveFileResult, error) {
	user, err := handler.userService.GetUser(ctx, args.SessionId)
	if err != nil {
		return MoveFileResult{}, err
	}

	drive, err := handler.driveRepo.GetDriveByOwner(ctx, user.Account())
	if err != nil {
		return MoveFileResult{}, err
	}

	file, err := handler.fileRepo.GetFile(ctx, args.FileId)
	if err != nil {
		return MoveFileResult{}, err
	}

	err = file.CanWrite(drive.Id())
	if err != nil {
		return MoveFileResult{}, err
	}

	if args.Parent.IsFolder() {
		parentFolder, err := handler.folderRepo.GetFolder(ctx, args.Parent.FolderId())
		if err != nil {
			// Move to folder that does not exist.
			return MoveFileResult{}, err
		}

		err = parentFolder.CanWrite(drive.Id())
		if err != nil {
			// Cannot write to parent folder.
			return MoveFileResult{}, err
		}
	}

	err = file.ChangeParent(args.Parent)
	if err != nil {
		return MoveFileResult{}, err
	}

	// Only once transaction.
	err = handler.fileRepo.UpdateFile(
		ctx,
		file,
		repository.UpdateFileOptions{
			MustInSameState: true,
		},
	)
	if err != nil {
		return MoveFileResult{}, err
	}

	return MoveFileResult{}, err
}

type MoveFileHandler decorator.Handler[MoveFileArgs, MoveFileResult]

func NewMoveFileHandler(
	userService remote.UserService,
	driveRepo repository.DriveRepository,
	folderRepo repository.FolderRepository,
	fileRepo repository.FileRepository,
	logger *logrus.Entry,
) MoveFileHandler {
	return decorator.WithLogging[MoveFileArgs, MoveFileResult](
		moveFileHandler{
			userService: userService,
			driveRepo:   driveRepo,
			folderRepo:  folderRepo,
			fileRepo:    fileRepo,
		},
		logger,
	)
}
