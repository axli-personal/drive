package domain

import (
	"errors"
	"io"
)

var (
	ErrInvalidFile = errors.New("invalid file")
	ErrInvalidData = errors.New("invalid data")
)

type Object struct {
	fileId     string
	totalBytes int
	data       io.Reader
}

func (object *Object) FileId() string {
	return object.fileId
}

func (object *Object) Read(p []byte) (int, error) {
	n, err := object.data.Read(p)

	object.totalBytes += n

	return n, err
}

func (object *Object) TotalBytes() int {
	return object.totalBytes
}

func NewObject(fileId string, data io.Reader) (*Object, error) {
	if fileId == "" {
		return nil, ErrInvalidFile
	}
	if data == nil {
		return nil, ErrInvalidData
	}

	return &Object{
		fileId: fileId,
		data:   data,
	}, nil
}
