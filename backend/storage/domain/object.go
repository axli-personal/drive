package domain

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"io"
)

var (
	ErrInvalidData = errors.New("invalid data")
)

type Object struct {
	hash       string
	totalBytes int
	data       io.Reader
}

func (object *Object) Hash() string {
	return object.hash
}

func (object *Object) Read(p []byte) (int, error) {
	return object.data.Read(p)
}

func (object *Object) TotalBytes() int {
	return object.totalBytes
}

func NewObject(data io.ReadSeeker) (*Object, error) {
	if data == nil {
		return nil, ErrInvalidData
	}

	sha1Hash := sha1.New()

	totalBytes, err := io.Copy(sha1Hash, data)
	if err != nil {
		return nil, err
	}

	_, err = data.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}

	hash := base64.RawURLEncoding.EncodeToString(sha1Hash.Sum(nil))

	return &Object{
		hash:       hash,
		totalBytes: int(totalBytes),
		data:       data,
	}, nil
}
