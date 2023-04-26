package usecases

import (
	"context"
	"github.com/axli-personal/drive/backend/common/decorator"
	"github.com/axli-personal/drive/backend/drive/domain"
	"github.com/axli-personal/drive/backend/drive/remote"
	"github.com/axli-personal/drive/backend/drive/repository"
	"github.com/axli-personal/drive/backend/pkg/errors"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const (
	ErrCodeFileNameExist = "FileNameExist"
)

type (
	StartUploadArgs struct {
		SessionId  string
		FileParent domain.Parent
		FileName   string
		FileHash   string
		FileSize   int
	}

	StartUploadResult struct {
		FileId uuid.UUID
	}
)

type startUploadHandler struct {
	repo        repository.Repository
	userService remote.UserService
}

func (h startUploadHandler) Handle(ctx context.Context, args StartUploadArgs) (result StartUploadResult, err error) {
	user, err := h.userService.GetUser(ctx, args.SessionId)
	if err != nil {
		return result, err
	}

	err = h.repo.Transaction(func(repo repository.Repository) error {
		drive, err := repo.GetDriveRepo().GetDriveByOwner(ctx, user.Account())
		if err != nil {
			return err
		}

		files, err := repo.GetFileRepo().FindFile(
			ctx,
			repository.FindFileOptions{
				DriveId: drive.Id(),
				Parent:  args.FileParent,
				Name:    args.FileName,
			},
		)
		if err != nil {
			return err
		}

		if len(files) > 0 {
			if files[0].State() != domain.StateLocked {
				return errors.New(ErrCodeFileNameExist, "duplicated file", nil)
			}

			result.FileId = files[0].Id()

			err = drive.IncreaseUsage(args.FileSize - files[0].Size())
			if err != nil {
				return err
			}

			files[0].SetHash(args.FileHash)
			files[0].SetSize(args.FileSize)

			err = repo.GetFileRepo().UpdateFile(ctx, files[0], repository.UpdateFileOptions{})
			if err != nil {
				return err
			}
		} else {
			file, err := domain.NewFile(drive.Id(), args.FileParent, args.FileName, args.FileSize, args.FileHash)
			if err != nil {
				return err
			}

			result.FileId = file.Id()

			err = drive.IncreaseUsage(args.FileSize)
			if err != nil {
				return err
			}

			err = repo.GetFileRepo().SaveFile(ctx, file)
			if err != nil {
				return err
			}
		}

		err = repo.GetDriveRepo().UpdateDrive(ctx, drive)
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}

type StartUploadHandler decorator.Handler[StartUploadArgs, StartUploadResult]

func NewStartUploadHandler(
	repo repository.Repository,
	userService remote.UserService,
	logger *logrus.Entry,
) StartUploadHandler {
	return decorator.WithLogging[StartUploadArgs, StartUploadResult](
		startUploadHandler{
			repo:        repo,
			userService: userService,
		},
		logger,
	)
}
