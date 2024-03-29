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
	RemoveFolderArgs struct {
		SessionId string
		FolderId  uuid.UUID
	}

	RemoveFolderResult struct {
	}
)

type removeFolderHandler struct {
	repo        repository.Repository
	userService remote.UserService
}

func (h removeFolderHandler) Handle(ctx context.Context, args RemoveFolderArgs) (RemoveFolderResult, error) {
	user, err := h.userService.GetUser(ctx, args.SessionId)
	if err != nil {
		return RemoveFolderResult{}, err
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

		err = folder.ManuallyTrash()
		if err != nil {
			return err
		}
		err = repo.GetFolderRepo().UpdateFolder(ctx, folder, repository.UpdateFolderOptions{UpdateChildrenState: true})
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
				err = subFiles[i].RecursivelyTrash()
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
				err = subFolders[i].RecursivelyTrash()
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
		return RemoveFolderResult{}, err
	}

	return RemoveFolderResult{}, nil
}

type RemoveFolderHandler decorator.Handler[RemoveFolderArgs, RemoveFolderResult]

func NewRemoveFolderHandler(
	repo repository.Repository,
	userService remote.UserService,
	logger *logrus.Entry,
) RemoveFolderHandler {
	return decorator.WithLogging[RemoveFolderArgs, RemoveFolderResult](
		removeFolderHandler{
			repo:        repo,
			userService: userService,
		},
		logger,
	)
}
