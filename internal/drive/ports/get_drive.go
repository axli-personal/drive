package ports

import (
	"github.com/axli-personal/drive/internal/common/auth"
	"github.com/axli-personal/drive/internal/drive/usecases"
	"github.com/axli-personal/drive/internal/pkg/types"
	"github.com/gofiber/fiber/v2"
)

func (server HTTPServer) GetDrive(ctx *fiber.Ctx) (err error) {
	sessionId := ctx.Cookies(auth.SessionIdCookieKey)
	if sessionId == "" {
		return ctx.SendStatus(fiber.StatusForbidden)
	}

	result, err := server.svc.GetDrive.Handle(
		ctx.Context(),
		usecases.GetDriveArgs{
			SessionId: sessionId,
		})
	if err != nil {
		return err
	}

	response := types.GetDriveResponse{
		DriveId:   result.Id.String(),
		Children:  types.Children{},
		PlanName:  result.Plan.Name(),
		UsedBytes: result.Usage.Bytes(),
		MaxBytes:  result.Plan.MaxBytes(),
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
