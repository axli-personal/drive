package usecases

import (
	"context"
	"errors"
	"github.com/axli-personal/drive/internal/common/decorator"
	"github.com/axli-personal/drive/internal/user/domain"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"time"
)

var (
	ErrWrongPassword = errors.New("wrong password")
)

type LoginArgs struct {
	Account  domain.Account
	Password domain.Password
}

type LoginResult struct {
	SessionId uuid.UUID
}

type loginHandler struct {
	userRepo    domain.UserRepository
	sessionRepo domain.SessionRepository
}

func (handler loginHandler) Handle(ctx context.Context, args LoginArgs) (LoginResult, error) {
	user, err := handler.userRepo.GetUser(ctx, args.Account)
	if err != nil {
		return LoginResult{}, err
	}

	if args.Password != user.Password() {
		return LoginResult{}, ErrWrongPassword
	}

	session, err := domain.NewSession(user.Account(), user.Username())
	if err != nil {
		return LoginResult{}, err
	}

	err = handler.sessionRepo.SaveSession(ctx, session, 12*time.Hour)
	if err != nil {
		return LoginResult{}, err
	}

	return LoginResult{
		SessionId: session.Id(),
	}, nil
}

type LoginHandler decorator.Handler[LoginArgs, LoginResult]

func NewLoginHandler(
	userRepo domain.UserRepository,
	sessionRepo domain.SessionRepository,
	logger *logrus.Entry,
) LoginHandler {
	return decorator.WithLogging[LoginArgs, LoginResult](
		loginHandler{
			sessionRepo: sessionRepo,
			userRepo:    userRepo,
		},
		logger,
	)
}
