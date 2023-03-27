package ports

import (
	"context"
	"github.com/axli-personal/drive/backend/drive/domain"
	"github.com/axli-personal/drive/backend/drive/usecases"
	"github.com/axli-personal/drive/backend/pkg/types"
)

func (server RPCServer) StartUpload(request *types.StartUploadRequest, response *types.StartUploadResponse) (err error) {
	parent, err := domain.CreateParent(request.FileParent)
	if err != nil {
		return err
	}

	result, err := server.svc.StartUpload.Handle(
		context.Background(),
		usecases.StartUploadArgs{
			SessionId:  request.SessionId,
			FileParent: parent,
			FileName:   request.FileName,
			FileHash:   request.FileHash,
			FileSize:   request.FileSize,
		},
	)
	if err != nil {
		return err
	}

	response.FileId = result.FileId.String()

	return nil
}
