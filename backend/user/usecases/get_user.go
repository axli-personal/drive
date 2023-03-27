package usecases

import (
	"context"
	"github.com/axli-personal/drive/backend/common/decorator"
	"github.com/axli-personal/drive/backend/pkg/errors"
	"github.com/axli-personal/drive/backend/user/domain"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var (
	ErrCodeUseCase  = "UseCase"
	ErrCodeNotLogin = "NotLogin"
)

type GetUserArgs struct {
	SessionId uuid.UUID
}

type GetUserResult struct {
	Account      domain.Account
	Username     string
	Introduction string
}

type getUserHandler struct {
	userRepo    domain.UserRepository
	sessionRepo domain.SessionRepository
}

func (handler getUserHandler) Handle(ctx context.Context, args GetUserArgs) (result GetUserResult, err error) {
	session, err := handler.sessionRepo.GetSession(ctx, args.SessionId)
	if err != nil {
		if err, ok := err.(*errors.Error); ok {
			if err.Code() == domain.ErrCodeNotFound {
				return GetUserResult{}, errors.New(ErrCodeNotLogin, "please login first", err)
			}
		}
		return GetUserResult{}, errors.New(ErrCodeUseCase, "fail to get user", err)
	}

	user, err := handler.userRepo.GetUser(ctx, session.Account())
	if err != nil {
		return GetUserResult{}, err
	}

	return GetUserResult{
		Account:      user.Account(),
		Username:     user.Username(),
		Introduction: user.Introduction(),
	}, nil
}

type GetUserHandler decorator.Handler[GetUserArgs, GetUserResult]

func NewGetUserHandler(
	userRepo domain.UserRepository,
	sessionRepo domain.SessionRepository,
	logger *logrus.Entry,
) GetUserHandler {
	return decorator.WithLogging[GetUserArgs, GetUserResult](
		getUserHandler{
			userRepo:    userRepo,
			sessionRepo: sessionRepo,
		},
		logger,
	)
}
