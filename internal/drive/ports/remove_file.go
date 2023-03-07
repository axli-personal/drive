package ports

import (
	"github.com/axli-personal/drive/internal/common/auth"
	"github.com/axli-personal/drive/internal/drive/usecases"
	"github.com/gofiber/fiber/v2"
)

func (server HTTPServer) RemoveFile(ctx *fiber.Ctx) (err error) {
	sessionId := ctx.Cookies(auth.SessionIdCookieKey)
	if sessionId == "" {
		return ctx.SendStatus(fiber.StatusForbidden)
	}

	_, err = server.svc.RemoveFile.Handle(
		ctx.Context(),
		usecases.RemoveFileArgs{
			SessionId: sessionId,
		},
	)

	return err
}
