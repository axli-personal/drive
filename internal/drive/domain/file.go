package domain

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

type File struct {
	id             uuid.UUID
	endpoint       string
	bytes          int
	downloadCounts int
	Metadata
}

var (
	ErrInvalidUUID                 = errors.New("invalid uuid")
	ErrCannotBindStorage           = errors.New("cannot bind storage")
	ErrCannotIncreaseDownloadCount = errors.New("cannot increase download count")
)

func NewFile(driveId uuid.UUID, parent Parent, fileName string) (*File, error) {
	if driveId == uuid.Nil {
		return nil, ErrInvalidUUID
	}
	if parent.IsZero() {
		return nil, ErrInvalidParent
	}

	return &File{
		id: uuid.New(),
		Metadata: Metadata{
			driveId:    driveId,
			parent:     parent,
			name:       fileName,
			state:      StateLocked,
			lastChange: time.Now(),
		},
		downloadCounts: 0,
	}, nil
}

func NewFileFromRepository(
	id uuid.UUID,
	endpoint string,
	bytes int,
	downloadCounts int,
	driveId uuid.UUID,
	parent Parent,
	name string,
	state State,
	lastChange time.Time,
) (*File, error) {
	return &File{
		id:             id,
		endpoint:       endpoint,
		bytes:          bytes,
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

func (file *File) Endpoint() string {
	return file.endpoint
}

func (file *File) Bytes() int {
	return file.bytes
}

func (file *File) DownloadCounts() int {
	return file.downloadCounts
}

func (file *File) BindStorage(endpoint string, bytes int) error {
	if file.Metadata.state != StateLocked {
		return ErrCannotBindStorage
	}

	file.Metadata.state = StatePrivate
	file.endpoint = endpoint
	file.bytes = bytes
	file.Metadata.lastChange = time.Now()

	return nil
}

func (file *File) IncreaseDownloadTimes() error {
	if file.Metadata.state != StateShared {
		return ErrCannotIncreaseDownloadCount
	}

	file.downloadCounts += 1

	return nil
}
