package ports

import (
	"context"
	"github.com/axli-personal/drive/backend/common/auth"
	"github.com/axli-personal/drive/backend/drive/usecases"
	"github.com/axli-personal/drive/backend/pkg/errors"
	"github.com/axli-personal/drive/backend/pkg/types"
	"github.com/gofiber/fiber/v2"
)

func (server HTTPServer) GetRecycleBin(ctx *fiber.Ctx) (err error) {
	sessionId := ctx.Cookies(auth.SessionIdCookieKey)
	if sessionId == "" {
		return ctx.Status(fiber.StatusForbidden).JSON(
			types.ErrorResponse{
				Code:    types.ErrCodeUnauthenticated,
				Message: "please login first",
				Detail:  "missing session cookie",
			},
		)
	}

	result, err := server.svc.GetRecycleBin.Handle(
		context.Background(),
		usecases.GetRecycleBinArgs{
			SessionId: sessionId,
		},
	)
	if err != nil {
		if err, ok := err.(*errors.Error); ok {
			errResponse := types.ErrorResponse{
				Code:    err.Code(),
				Message: err.Message(),
				Detail:  err.Error(),
			}
			if err.Code() == types.ErrCodeUnauthenticated {
				return ctx.Status(fiber.StatusForbidden).JSON(errResponse)
			}
			if err.Code() == usecases.ErrCodeNotCreateDrive {
				return ctx.Status(fiber.StatusNotFound).JSON(errResponse)
			}
			return ctx.Status(fiber.StatusInternalServerError).JSON(errResponse)
		}
		return err
	}

	response := types.GetRecycleBinResponse{
		Children: types.Children{},
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
