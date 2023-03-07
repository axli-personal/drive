package ports

import (
	"github.com/axli-personal/drive/internal/common/auth"
	"github.com/axli-personal/drive/internal/drive/domain"
	"github.com/axli-personal/drive/internal/drive/usecases"
	"github.com/gofiber/fiber/v2"
)

func (server HTTPServer) GetPath(ctx *fiber.Ctx) (err error) {
	sessionId := ctx.Cookies(auth.SessionIdCookieKey)
	if sessionId == "" {
		return ctx.SendStatus(fiber.StatusForbidden)
	}

	parent, err := domain.CreateParent(ctx.Params("parent"))
	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	result, err := server.svc.GetPath.Handle(
		ctx.Context(),
		usecases.GetPathArgs{
			SessionId: sessionId,
			Parent:    parent,
		},
	)
	if err != nil {
		return err
	}

	return ctx.JSON(result)
}
