package ports

import (
	"github.com/axli-personal/drive/internal/common/auth"
	"github.com/axli-personal/drive/internal/pkg/types"
	"github.com/axli-personal/drive/internal/storage/usecases"
	"github.com/gofiber/fiber/v2"
	"path/filepath"
)

func (server HTTPServer) GetObject(ctx *fiber.Ctx) (err error) {
	sessionId := ctx.Cookies(auth.SessionIdCookieKey)
	if sessionId == "" {
		return ctx.SendStatus(fiber.StatusForbidden)
	}

	request := types.GetObjectRequest{}

	err = ctx.ParamsParser(&request)
	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	err = ctx.QueryParser(&request)
	if err != nil {
		request.Download = false
	}

	result, err := server.svc.DownloadObject.Handle(
		ctx.Context(),
		usecases.DownloadObjectArgs{
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
