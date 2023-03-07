package adapters

import (
	"github.com/axli-personal/drive/internal/pkg/types"
	"github.com/axli-personal/drive/internal/storage/remote"
	"net/rpc"
)

type RPCDriveService struct {
	client *rpc.Client
}

func NewRPCDriveService(address string) (remote.DriveService, error) {
	client, err := rpc.DialHTTP("tcp", address)

	return RPCDriveService{client: client}, err
}

func (service RPCDriveService) CanUpload(request types.CanUploadRequest) (types.CanUploadResponse, error) {
	response := types.CanUploadResponse{}

	err := service.client.Call("RPCServer.CanUpload", &request, &response)

	return response, err
}

func (service RPCDriveService) CanDownload(request types.CanDownloadRequest) (types.CanDownloadResponse, error) {
	response := types.CanDownloadResponse{}

	err := service.client.Call("RPCServer.CanDownload", &request, &response)

	return response, err
}
