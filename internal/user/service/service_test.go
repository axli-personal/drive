package service_test

import (
	"context"
	"github.com/axli-personal/drive/internal/common/utils"
	"github.com/axli-personal/drive/internal/user/domain"
	"github.com/axli-personal/drive/internal/user/service"
	"github.com/axli-personal/drive/internal/user/usecases"
	"github.com/caarlos0/env/v7"
	"testing"
)

func TestService(t *testing.T) {
	config := service.Config{}

	err := env.Parse(&config)
	if err != nil {
		t.Fatal(err)
	}

	svc, err := service.NewService(config)
	if err != nil {
		t.Fatal(err)
	}

	account, err := domain.NewAccount(utils.RandomString(10))
	password, err := domain.NewPassword(utils.RandomString(10))
	username := utils.RandomString(10)

	_, err = svc.Register.Handle(
		context.Background(),
		usecases.RegisterArgs{
			Account:  account,
			Password: password,
			Username: username,
		},
	)
	if err != nil {
		t.Error(err)
	}

	loginResult, err := svc.Login.Handle(
		context.Background(),
		usecases.LoginArgs{
			Account:  account,
			Password: password,
		},
	)
	if err != nil {
		t.Error(err)
	}

	_, err = svc.GetUser.Handle(
		context.Background(),
		usecases.GetUserArgs{
			SessionId: loginResult.SessionId,
		},
	)
	if err != nil {
		t.Error(err)
	}
}
