package ports

import (
	"github.com/axli-personal/drive/internal/common/auth"
	"github.com/axli-personal/drive/internal/drive/usecases"
	"github.com/axli-personal/drive/internal/pkg/types"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (server HTTPServer) GetFile(ctx *fiber.Ctx) (err error) {
	sessionId := ctx.Cookies(auth.SessionIdCookieKey)
	if sessionId == "" {
		return ctx.SendStatus(fiber.StatusForbidden)
	}

	request := types.GetFileRequest{}

	err = ctx.ParamsParser(&request)
	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	fileId, err := uuid.Parse(request.FileId)
	if err != nil || fileId == uuid.Nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	result, err := server.svc.GetFile.Handle(
		ctx.Context(),
		usecases.GetFileArgs{
			SessionId: sessionId,
			FileId:    fileId,
		},
	)
	if err != nil {
		return err
	}

	response := types.GetFileResponse{
		FileId:     result.FileId.String(),
		Parent:     result.Parent.String(),
		Name:       result.Name,
		LastChange: result.LastChange,
		Endpoint:   result.Endpoint,
		Bytes:      result.Bytes,
	}

	return ctx.JSON(response)
}
