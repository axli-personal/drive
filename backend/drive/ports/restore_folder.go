package ports

import (
	"context"
	"github.com/axli-personal/drive/backend/common/auth"
	"github.com/axli-personal/drive/backend/drive/usecases"
	"github.com/axli-personal/drive/backend/pkg/types"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (server HTTPServer) RestoreFolder(ctx *fiber.Ctx) (err error) {
	sessionId := ctx.Cookies(auth.SessionIdCookieKey)
	if sessionId == "" {
		return ctx.SendStatus(fiber.StatusForbidden)
	}

	request := types.RestoreFolderRequest{}

	err = ctx.ParamsParser(&request)
	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	folderId, err := uuid.Parse(request.FolderId)
	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	_, err = server.svc.RestoreFolder.Handle(
		context.Background(),
		usecases.RestoreFolderArgs{
			SessionId: sessionId,
			FolderId:  folderId,
		},
	)

	return err
}
