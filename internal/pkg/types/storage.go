package types

// HTTP

type GetObjectRequest struct {
	FileId   string `params:"fileId"`
	Download bool   `query:"download"`
}
