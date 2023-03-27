package ports

import (
	"github.com/axli-personal/drive/backend/pkg/types"
	"github.com/axli-personal/drive/backend/user/domain"
	"github.com/axli-personal/drive/backend/user/usecases"
	"github.com/gofiber/fiber/v2"
)

func (server HTTPServer) Register(ctx *fiber.Ctx) (err error) {
	request := types.RegisterRequest{}

	err = ctx.BodyParser(&request)
	if err != nil {
		return err
	}

	account, err := domain.NewAccount(request.Account)
	if err != nil {
		return err
	}
	password, err := domain.NewPassword(request.Password)
	if err != nil {
		return err
	}

	_, err = server.svc.Register.Handle(
		ctx.Context(),
		usecases.RegisterArgs{
			Account:  account,
			Password: password,
			Username: request.Username,
		},
	)

	return err
}
