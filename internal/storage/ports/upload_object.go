package ports

import (
	"github.com/axli-personal/drive/internal/common/auth"
	"github.com/axli-personal/drive/internal/storage/usecases"
	"github.com/gofiber/fiber/v2"
)

func (server HTTPServer) UploadObject(ctx *fiber.Ctx) (err error) {
	sessionId := ctx.Cookies(auth.SessionIdCookieKey)
	if sessionId == "" {
		return ctx.SendStatus(fiber.StatusForbidden)
	}

	fileId := ctx.Params("fileId")
	if fileId == "" {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		return err
	}

	files := form.File["files"]

	if len(files) != 1 {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	data, err := files[0].Open()
	if err != nil {
		return err
	}

	_, err = server.svc.UploadObject.Handle(
		ctx.Context(),
		usecases.UploadObjectArgs{
			SessionId: sessionId,
			FileId:    fileId,
			Data:      data,
		},
	)
	if err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusOK)
}
