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

type (
	RestoreFolderArgs struct {
		SessionId string
		FolderId  uuid.UUID
	}

	RestoreFolderResult struct {
	}
)

type restoreFolderHandler struct {
	repo        repository.Repository
	userService remote.UserService
}

func (h restoreFolderHandler) Handle(ctx context.Context, args RestoreFolderArgs) (RestoreFolderResult, error) {
	user, err := h.userService.GetUser(ctx, args.SessionId)
	if err != nil {
		return RestoreFolderResult{}, err
	}

	err = h.repo.Transaction(func(repo repository.Repository) error {
		drive, err := repo.GetDriveRepo().GetDriveByOwner(ctx, user.Account())
		if err != nil {
			return err
		}

		folder, err := h.repo.GetFolderRepo().GetFolder(ctx, args.FolderId)
		if err != nil {
			return err
		}

		err = folder.CanWrite(drive.Id())
		if err != nil {
			return err
		}

		err = folder.ManuallyRestore()
		if err != nil {
			return err
		}
		err = repo.GetFolderRepo().UpdateFolder(ctx, folder, repository.UpdateFolderOptions{})
		if err != nil {
			return err
		}

		folders := []*domain.Folder{folder}
		for len(folders) > 0 {
			parent, err := domain.CreateFolderParent(folders[0].Id())
			if err != nil {
				return err
			}

			subFiles, err := repo.GetFileRepo().FindFile(
				ctx,
				repository.FindFileOptions{
					Parent: parent,
				},
			)
			if err != nil {
				return err
			}
			for i := 0; i < len(subFiles); i++ {
				err = subFiles[i].RecursivelyRestore()
				if err != nil {
					return err
				}
				err = repo.GetFileRepo().UpdateFile(ctx, subFiles[i], repository.UpdateFileOptions{})
				if err != nil {
					return err
				}
			}

			subFolders, err := repo.GetFolderRepo().FindFolder(
				ctx,
				repository.FindFolderOptions{
					Parent: parent,
				},
			)
			if err != nil {
				return err
			}
			for i := 0; i < len(subFolders); i++ {
				err = subFolders[i].RecursivelyRestore()
				if err != nil {
					return err
				}
				err = repo.GetFolderRepo().UpdateFolder(ctx, subFolders[i], repository.UpdateFolderOptions{})
				if err != nil {
					return err
				}
			}

			folders = append(folders[1:], subFolders...)
		}

		return nil
	})
	if err != nil {
		return RestoreFolderResult{}, err
	}

	return RestoreFolderResult{}, nil
}

type RestoreFolderHandler decorator.Handler[RestoreFolderArgs, RestoreFolderResult]

func NewRestoreFolderHandler(
	repo repository.Repository,
	userService remote.UserService,
	logger *logrus.Entry,
) RestoreFolderHandler {
	return decorator.WithLogging[RestoreFolderArgs, RestoreFolderResult](
		restoreFolderHandler{
			repo:        repo,
			userService: userService,
		},
		logger,
	)
}
