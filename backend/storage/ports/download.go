package ports

import (
	"github.com/axli-personal/drive/backend/common/auth"
	"github.com/axli-personal/drive/backend/pkg/types"
	"github.com/axli-personal/drive/backend/storage/usecases"
	"github.com/gofiber/fiber/v2"
	"path/filepath"
)

func (server HTTPServer) Download(ctx *fiber.Ctx) (err error) {
	sessionId := ctx.Cookies(auth.SessionIdCookieKey)

	request := types.GetObjectRequest{}

	err = ctx.ParamsParser(&request)
	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	err = ctx.QueryParser(&request)
	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	result, err := server.svc.DownloadObject.Handle(
		ctx.Context(),
		usecases.DownloadArgs{
			SessionId: sessionId,
			FileId:    request.FileId,
		},
	)
	if err != nil {
		return err
	}

	if request.Download {
		ctx.Attachment(result.FileName)
	}

	return ctx.Type(filepath.Ext(result.FileName)).SendStream(result.Data)
}
