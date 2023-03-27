package ports

import (
	"github.com/axli-personal/drive/backend/common/auth"
	"github.com/axli-personal/drive/backend/drive/usecases"
	"github.com/axli-personal/drive/backend/pkg/types"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (server HTTPServer) GetFolder(ctx *fiber.Ctx) (err error) {
	sessionId := ctx.Cookies(auth.SessionIdCookieKey)

	request := types.GetFolderRequest{}

	err = ctx.ParamsParser(&request)
	if err != nil {
		return err
	}

	folderId, err := uuid.Parse(request.FolderId)
	if err != nil || folderId == uuid.Nil {
		return err
	}

	result, err := server.svc.GetFolder.Handle(
		ctx.Context(),
		usecases.GetFolderArgs{
			SessionId: sessionId,
			FolderId:  folderId,
		},
	)
	if err != nil {
		return err
	}

	response := types.GetFolderResponse{
		FolderId:   result.FolderId.String(),
		Parent:     result.Parent.String(),
		Name:       result.Name,
		Shared:     result.Shared,
		LastChange: result.LastChange,
		Children:   types.Children{},
	}

	for _, link := range result.Children.Folders {
		response.Children.Folders = append(response.Children.Folders, types.FolderLink{
			Id:   link.Id.String(),
			Name: link.Name,
		})
	}

	for _, link := range result.Children.Files {
		response.Children.Files = append(response.Children.Files, types.FileLink{
			Id:    link.Id.String(),
			Name:  link.Name,
			Bytes: link.Bytes,
		})
	}

	return ctx.JSON(response)
}
