package ports

import (
	"context"
	"github.com/axli-personal/drive/internal/drive/usecases"
	"github.com/axli-personal/drive/internal/pkg/types"
	"github.com/google/uuid"
)

func (server RPCServer) CanUpload(request *types.CanUploadRequest, response *types.CanUploadResponse) (err error) {
	sessionId, err := uuid.Parse(request.SessionId)
	if err != nil {
		return err
	}

	fileId, err := uuid.Parse(request.FileId)
	if err != nil {
		return err
	}

	result, err := server.svc.CanUpload.Handle(
		context.Background(),
		usecases.CanUploadArgs{
			SessionId: sessionId.String(),
			FileId:    fileId,
		},
	)
	if err != nil {
		return err
	}

	response.CanUpload = result.CanUpload

	return nil
}
