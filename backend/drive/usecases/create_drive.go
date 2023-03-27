package usecases

import (
	"context"
	"github.com/axli-personal/drive/backend/common/decorator"
	"github.com/axli-personal/drive/backend/drive/domain"
	"github.com/axli-personal/drive/backend/drive/remote"
	"github.com/axli-personal/drive/backend/drive/repository"
	"github.com/sirupsen/logrus"
)

type (
	CreateDriveArgs struct {
		SessionId string
	}

	CreateDriveResult struct {
	}
)

type createDriveHandler struct {
	userService remote.UserService
	driveRepo   repository.DriveRepository
}

func (handler createDriveHandler) Handle(ctx context.Context, args CreateDriveArgs) (CreateDriveResult, error) {
	user, err := handler.userService.GetUser(ctx, args.SessionId)
	if err != nil {
		return CreateDriveResult{}, err
	}

	drive, err := domain.NewDrive(user.Account())
	if err != nil {
		return CreateDriveResult{}, err
	}

	err = handler.driveRepo.CreateDrive(
		ctx,
		drive,
		repository.CreateDriveOptions{
			OnlyOneDrive: true,
		},
	)
	if err != nil {
		return CreateDriveResult{}, err
	}

	return CreateDriveResult{}, nil
}

type CreateDriveHandler decorator.Handler[CreateDriveArgs, CreateDriveResult]

func NewCreateDriveHandler(
	userService remote.UserService,
	driveRepo repository.DriveRepository,
	logger *logrus.Entry,
) CreateDriveHandler {
	return decorator.WithLogging[CreateDriveArgs, CreateDriveResult](
		createDriveHandler{
			userService: userService,
			driveRepo:   driveRepo,
		},
		logger,
	)
}
