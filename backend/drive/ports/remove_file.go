package ports

import (
	"context"
	"github.com/axli-personal/drive/backend/common/auth"
	"github.com/axli-personal/drive/backend/drive/usecases"
	"github.com/axli-personal/drive/backend/pkg/types"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (server HTTPServer) RemoveFile(ctx *fiber.Ctx) (err error) {
	sessionId := ctx.Cookies(auth.SessionIdCookieKey)
	if sessionId == "" {
		return ctx.SendStatus(fiber.StatusForbidden)
	}

	request := types.RemoveFileRequest{}

	err = ctx.ParamsParser(&request)
	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	fileId, err := uuid.Parse(request.FileId)
	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	_, err = server.svc.RemoveFile.Handle(
		context.Background(),
		usecases.RemoveFileArgs{
			SessionId: sessionId,
			FileId:    fileId,
		},
	)

	return err
}
