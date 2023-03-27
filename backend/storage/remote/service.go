package remote

import (
	"github.com/axli-personal/drive/backend/pkg/types"
)

type (
	DriveService interface {
		StartUpload(request types.StartUploadRequest) (types.StartUploadResponse, error)
		StartDownload(request types.StartDownloadRequest) (types.StartDownloadResponse, error)
		FinishUpload(request types.FinishUploadRequest) (types.FinishUploadResponse, error)
	}
)
