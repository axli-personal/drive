package types

import "time"

// RPC

type CanDownloadRequest struct {
	SessionId string
	FileId    string
}

type CanDownloadResponse struct {
	CanDownload bool
	FileName    string
}

type CanUploadRequest struct {
	SessionId string
	FileId    string
}

type CanUploadResponse struct {
	CanUpload bool
}

// HTTP

type CreateFileRequest struct {
	Parent   string `json:"parent"`
	FileName string `json:"fileName"`
}

type CreateFileResponse struct {
	FileId          string    `json:"fileId"`
	FileName        string    `json:"fileName"`
	LastChange      time.Time `json:"lastChange"`
	StorageEndpoint string    `json:"storageEndpoint"`
}

type CreateFolderRequest struct {
	Parent     string `json:"parent"`
	FolderName string `json:"folderName"`
}

type GetFileRequest struct {
	FileId string `params:"fileId"`
}

type GetFileResponse struct {
	FileId         string    `json:"fileId"`
	Parent         string    `json:"parent"`
	Name           string    `json:"name"`
	LastChange     time.Time `json:"lastChange"`
	Endpoint       string    `json:"endpoint"`
	Bytes          int       `json:"bytes"`
	DownloadCounts int       `json:"downloadCounts"`
}

type GetFolderRequest struct {
	FolderId string `params:"folderId"`
}

type GetFolderResponse struct {
	FolderId   string    `json:"folderId"`
	Parent     string    `json:"parent"`
	Name       string    `json:"name"`
	LastChange time.Time `json:"lastChange"`
	Children   Children  `json:"children"`
}

type Children struct {
	Folders []FolderLink `json:"folders"`
	Files   []FileLink   `json:"files"`
}

type FolderLink struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type FileLink struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Bytes int    `json:"bytes"`
}

type GetDriveResponse struct {
	DriveId   string   `json:"driveId"`
	Children  Children `json:"children"`
	PlanName  string   `json:"planName"`
	UsedBytes int      `json:"usedBytes"`
	MaxBytes  int      `json:"maxBytes"`
}
