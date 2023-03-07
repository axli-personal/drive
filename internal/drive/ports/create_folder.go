package ports

import (
	"github.com/axli-personal/drive/internal/common/auth"
	"github.com/axli-personal/drive/internal/drive/domain"
	"github.com/axli-personal/drive/internal/drive/usecases"
	"github.com/axli-personal/drive/internal/pkg/types"
	"github.com/gofiber/fiber/v2"
)

func (server HTTPServer) CreateFolder(ctx *fiber.Ctx) (err error) {
	request := types.CreateFolderRequest{}

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

	_, err = server.svc.CreateFolder.Handle(
		ctx.Context(),
		usecases.CreateFolderArgs{
			SessionId:    sessionId,
			FolderParent: parent,
			FolderName:   request.FolderName,
		},
	)
	if err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusOK)
}
