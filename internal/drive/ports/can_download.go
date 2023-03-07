package ports

import (
	"context"
	"github.com/axli-personal/drive/internal/drive/usecases"
	"github.com/axli-personal/drive/internal/pkg/types"
	"github.com/google/uuid"
)

func (server RPCServer) CanDownload(request *types.CanDownloadRequest, response *types.CanDownloadResponse) (err error) {
	sessionId, err := uuid.Parse(request.SessionId)
	if err != nil {
		return err
	}

	fileId, err := uuid.Parse(request.FileId)
	if err != nil {
		return err
	}

	result, err := server.svc.CanDownload.Handle(
		context.Background(),
		usecases.CanDownloadArgs{
			SessionId: sessionId.String(),
			FileId:    fileId,
		},
	)
	if err != nil {
		return err
	}

	response.CanDownload = result.CanDownload
	response.FileName = result.FileName

	return nil
}
