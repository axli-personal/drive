package ports

import (
	"github.com/axli-personal/drive/backend/common/auth"
	"github.com/axli-personal/drive/backend/storage/usecases"
	"github.com/gofiber/fiber/v2"
)

func (server HTTPServer) Upload(ctx *fiber.Ctx) (err error) {
	sessionId := ctx.Cookies(auth.SessionIdCookieKey)
	if sessionId == "" {
		return ctx.SendStatus(fiber.StatusForbidden)
	}

	fileParent := ctx.FormValue("parent")
	if fileParent == "" {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	file, err := fileHeader.Open()
	if err != nil {
		return err
	}

	_, err = server.svc.UploadObject.Handle(
		ctx.Context(),
		usecases.UploadArgs{
			SessionId:  sessionId,
			FileParent: fileParent,
			FileName:   fileHeader.Filename,
			Data:       file,
		},
	)
	if err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusOK)
}
