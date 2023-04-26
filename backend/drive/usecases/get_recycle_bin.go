package usecases

import (
	"context"
	"github.com/axli-personal/drive/backend/common/decorator"
	"github.com/axli-personal/drive/backend/drive/domain"
	"github.com/axli-personal/drive/backend/drive/remote"
	"github.com/axli-personal/drive/backend/drive/repository"
	"github.com/axli-personal/drive/backend/pkg/errors"
	"github.com/sirupsen/logrus"
)

type (
	GetRecycleBinArgs struct {
		SessionId string
	}

	GetRecycleBinResult struct {
		Children Children
	}
)

type getRecycleBinHandler struct {
	repo        repository.Repository
	userService remote.UserService
}

func (h getRecycleBinHandler) Handle(ctx context.Context, args GetRecycleBinArgs) (GetRecycleBinResult, error) {
	user, err := h.userService.GetUser(ctx, args.SessionId)
	if err != nil {
		return GetRecycleBinResult{}, err
	}

	drive, err := h.repo.GetDriveRepo().GetDriveByOwner(ctx, user.Account())
	if err != nil {
		if err, ok := err.(*errors.Error); ok {
			if err.Code() == repository.ErrCodeNotFound {
				return GetRecycleBinResult{}, errors.New(ErrCodeNotCreateDrive, "please create drive first", err)
			}
		}
		return GetRecycleBinResult{}, errors.New(ErrCodeUseCase, "fail to get drive", err)
	}

	folders, err := h.repo.GetFolderRepo().FindFolder(
		ctx,
		repository.FindFolderOptions{
			DriveId: drive.Id(),
			States:  []domain.State{domain.StateTrashedRoot},
		},
	)

	files, err := h.repo.GetFileRepo().FindFile(
		ctx,
		repository.FindFileOptions{
			DriveId: drive.Id(),
			States:  []domain.State{domain.StateTrashedRoot},
		},
	)

	return GetRecycleBinResult{
		Children: ToChildren(folders, files),
	}, nil
}

type GetRecycleBinHandler decorator.Handler[GetRecycleBinArgs, GetRecycleBinResult]

func NewGetRecycleBinHandler(
	repo repository.Repository,
	userService remote.UserService,
	logger *logrus.Entry,
) GetRecycleBinHandler {
	return decorator.WithLogging[GetRecycleBinArgs, GetRecycleBinResult](
		getRecycleBinHandler{
			repo:        repo,
			userService: userService,
		},
		logger,
	)
}
