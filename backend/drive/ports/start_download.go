package ports

import (
	"context"
	"github.com/axli-personal/drive/backend/drive/usecases"
	"github.com/axli-personal/drive/backend/pkg/types"
	"github.com/google/uuid"
)

func (server RPCServer) StartDownload(request *types.StartDownloadRequest, response *types.StartDownloadResponse) (err error) {
	fileId, err := uuid.Parse(request.FileId)
	if err != nil {
		return err
	}

	result, err := server.svc.StartDownload.Handle(
		context.Background(),
		usecases.StartDownloadArgs{
			SessionId: request.SessionId,
			FileId:    fileId,
		},
	)
	if err != nil {
		return err
	}

	response.FileName = result.FileName
	response.FileHash = result.FileHash

	return nil
}
