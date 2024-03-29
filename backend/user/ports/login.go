package ports

import (
	"github.com/axli-personal/drive/backend/common/auth"
	"github.com/axli-personal/drive/backend/pkg/types"
	"github.com/axli-personal/drive/backend/user/domain"
	"github.com/axli-personal/drive/backend/user/usecases"
	"github.com/gofiber/fiber/v2"
)

func (server HTTPServer) Login(ctx *fiber.Ctx) (err error) {
	request := types.LoginRequest{}

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

	result, err := server.svc.Login.Handle(
		ctx.Context(),
		usecases.LoginArgs{
			Account:  account,
			Password: password,
		},
	)
	if err != nil {
		return err
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     auth.SessionIdCookieKey,
		Value:    result.SessionId.String(),
		SameSite: fiber.CookieSameSiteNoneMode,
		Secure:   true,
	})

	return nil
}
