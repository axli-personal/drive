package ports

import (
	"github.com/axli-personal/drive/backend/common/auth"
	"github.com/axli-personal/drive/backend/drive/usecases"
	"github.com/gofiber/fiber/v2"
)

func (server HTTPServer) CreateDrive(ctx *fiber.Ctx) (err error) {
	sessionId := ctx.Cookies(auth.SessionIdCookieKey)
	if sessionId == "" {
		return ctx.SendStatus(fiber.StatusForbidden)
	}

	_, err = server.svc.CreateDrive.Handle(
		ctx.Context(),
		usecases.CreateDriveArgs{
			SessionId: sessionId,
		},
	)
	if err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusOK)
}
