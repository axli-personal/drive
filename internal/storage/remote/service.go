package remote

import (
	"github.com/axli-personal/drive/internal/pkg/types"
)

type (
	DriveService interface {
		CanUpload(request types.CanUploadRequest) (types.CanUploadResponse, error)
		CanDownload(request types.CanDownloadRequest) (types.CanDownloadResponse, error)
	}
)
