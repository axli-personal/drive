package events

var (
	FieldBody            = "Body"
	StreamFileUploaded   = "file-uploaded"
	StreamFileDownloaded = "file-downloaded"
	StreamFileDeleted    = "file-deleted"
	StreamFolderRemoved  = "folder-removed"
)

type FileUploaded struct {
	EventId    string `json:"-"`
	Endpoint   string `json:"endpoint"`
	FileId     string `json:"fileId"`
	TotalBytes int    `json:"totalBytes"`
}

type FileDownloaded struct {
	EventId string `json:"-"`
	FileId  string `json:"fileId"`
}

type FileDeleted struct {
	EventId string `json:"-"`
	FileId  string `json:"fileId"`
}

type FolderRemoved struct {
	EventId  string `json:"-"`
	FolderId string `json:"folderId"`
}
