package usecases

import (
	"context"
	"github.com/axli-personal/drive/internal/common/decorator"
	"github.com/axli-personal/drive/internal/user/domain"
	"github.com/sirupsen/logrus"
)

type RegisterArgs struct {
	Account  domain.Account
	Password domain.Password
	Username string
}

type RegisterResult struct {
}

type registerHandler struct {
	userRepo domain.UserRepository
}

func (handler registerHandler) Handle(ctx context.Context, args RegisterArgs) (result RegisterResult, err error) {
	user, err := domain.NewUser(args.Account, args.Password, args.Username)
	if err != nil {
		return RegisterResult{}, err
	}

	err = handler.userRepo.SaveUser(ctx, user)
	if err != nil {
		return RegisterResult{}, err
	}

	return RegisterResult{}, nil
}

type RegisterHandler decorator.Handler[RegisterArgs, RegisterResult]

func NewRegisterHandler(
	userRepo domain.UserRepository,
	logger *logrus.Entry,
) RegisterHandler {
	return decorator.WithLogging[RegisterArgs, RegisterResult](
		registerHandler{
			userRepo: userRepo,
		},
		logger,
	)
}
