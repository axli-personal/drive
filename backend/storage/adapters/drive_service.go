package adapters

import (
	"github.com/axli-personal/drive/backend/pkg/types"
	"github.com/axli-personal/drive/backend/storage/remote"
	"net/rpc"
)

type RPCDriveService struct {
	client *rpc.Client
}

func NewRPCDriveService(address string) (remote.DriveService, error) {
	client, err := rpc.DialHTTP("tcp", address)

	return RPCDriveService{client: client}, err
}

func (service RPCDriveService) StartUpload(request types.StartUploadRequest) (types.StartUploadResponse, error) {
	response := types.StartUploadResponse{}

	err := service.client.Call("RPCServer.StartUpload", &request, &response)

	return response, err
}

func (service RPCDriveService) StartDownload(request types.StartDownloadRequest) (types.StartDownloadResponse, error) {
	response := types.StartDownloadResponse{}

	err := service.client.Call("RPCServer.StartDownload", &request, &response)

	return response, err
}

func (service RPCDriveService) FinishUpload(request types.FinishUploadRequest) (types.FinishUploadResponse, error) {
	response := types.FinishUploadResponse{}

	err := service.client.Call("RPCServer.FinishUpload", &request, &response)

	return response, err
}
