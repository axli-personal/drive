package ports

import (
	"context"
	"github.com/axli-personal/drive/backend/drive/usecases"
	"github.com/axli-personal/drive/backend/pkg/types"
	"github.com/google/uuid"
)

func (server RPCServer) FinishUpload(request *types.FinishUploadRequest, response *types.FinishUploadResponse) (err error) {
	fileId, err := uuid.Parse(request.FileId)
	if err != nil {
		return err
	}

	_, err = server.svc.FinishUpload.Handle(
		context.Background(),
		usecases.FinishUploadArgs{
			FileId: fileId,
		},
	)
	if err != nil {
		return err
	}

	return nil
}
