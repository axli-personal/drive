package usecases

import (
	"context"
	"github.com/axli-personal/drive/internal/common/decorator"
	"github.com/axli-personal/drive/internal/user/domain"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
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
		return GetUserResult{}, err
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
