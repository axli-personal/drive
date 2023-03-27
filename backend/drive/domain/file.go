package domain

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

type File struct {
	id             uuid.UUID
	size           int
	hash           string
	downloadCounts int
	Metadata
}

var (
	ErrInvalidUUID                 = errors.New("invalid uuid")
	ErrCannotIncreaseDownloadCount = errors.New("cannot increase download count")
)

func NewFile(driveId uuid.UUID, fileParent Parent, fileName string, fileSize int, fileHash string) (*File, error) {
	if driveId == uuid.Nil {
		return nil, ErrInvalidUUID
	}
	if fileParent.IsZero() {
		return nil, ErrInvalidParent
	}

	return &File{
		id: uuid.New(),
		Metadata: Metadata{
			driveId:    driveId,
			parent:     fileParent,
			name:       fileName,
			state:      StateLocked,
			lastChange: time.Now(),
		},
		size:           fileSize,
		hash:           fileHash,
		downloadCounts: 0,
	}, nil
}

func NewFileFromRepository(
	id uuid.UUID,
	size int,
	hash string,
	downloadCounts int,
	driveId uuid.UUID,
	parent Parent,
	name string,
	state State,
	lastChange time.Time,
) (*File, error) {
	return &File{
		id:             id,
		size:           size,
		hash:           hash,
		downloadCounts: downloadCounts,
		Metadata: Metadata{
			driveId:    driveId,
			parent:     parent,
			name:       name,
			state:      state,
			lastChange: lastChange,
		},
	}, nil
}

func (file *File) Id() uuid.UUID {
	return file.id
}

func (file *File) Size() int {
	return file.size
}

func (file *File) Hash() string {
	return file.hash
}

func (file *File) DownloadCounts() int {
	return file.downloadCounts
}

func (file *File) IncreaseDownloadCounts() {
	file.downloadCounts += 1
}
