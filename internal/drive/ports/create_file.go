package ports

import (
	"github.com/axli-personal/drive/internal/common/auth"
	"github.com/axli-personal/drive/internal/drive/domain"
	"github.com/axli-personal/drive/internal/drive/usecases"
	"github.com/axli-personal/drive/internal/pkg/types"
	"github.com/gofiber/fiber/v2"
)

func (server HTTPServer) CreateFile(ctx *fiber.Ctx) (err error) {
	request := types.CreateFileRequest{}

	err = ctx.BodyParser(&request)
	if err != nil {
		return err
	}

	sessionId := ctx.Cookies(auth.SessionIdCookieKey)
	if sessionId == "" {
		return ctx.SendStatus(fiber.StatusForbidden)
	}

	parent, err := domain.CreateParent(request.Parent)
	if err != nil {
		return err
	}

	result, err := server.svc.CreateFile.Handle(
		ctx.Context(),
		usecases.CreateFileArgs{
			SessionId:  sessionId,
			FileParent: parent,
			FileName:   request.FileName,
		},
	)
	if err != nil {
		return err
	}

	return ctx.JSON(types.CreateFileResponse{
		FileId:          result.FileId.String(),
		FileName:        result.FileName,
		LastChange:      result.LastChange,
		StorageEndpoint: result.StorageEndpoint,
	})
}
