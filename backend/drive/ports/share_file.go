package ports

import (
	"context"
	"github.com/axli-personal/drive/backend/common/auth"
	"github.com/axli-personal/drive/backend/drive/usecases"
	"github.com/axli-personal/drive/backend/pkg/types"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (server HTTPServer) ShareFile(ctx *fiber.Ctx) (err error) {
	sessionId := ctx.Cookies(auth.SessionIdCookieKey)
	if sessionId == "" {
		return ctx.SendStatus(fiber.StatusForbidden)
	}

	request := types.ShareFileRequest{}

	err = ctx.ParamsParser(&request)
	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	fileId, err := uuid.Parse(request.FileId)
	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	_, err = server.svc.ShareFile.Handle(
		context.Background(),
		usecases.ShareFileArgs{
			SessionId: sessionId,
			FileId:    fileId,
		},
	)

	return err
}
