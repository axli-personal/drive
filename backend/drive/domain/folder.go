package domain

import (
	"github.com/google/uuid"
	"time"
)

type Folder struct {
	id uuid.UUID
	Metadata
}

func NewFolder(driveId uuid.UUID, parent Parent, name string) (*Folder, error) {
	if driveId == uuid.Nil {
		return nil, ErrInvalidUUID
	}
	if parent.IsZero() {
		return nil, ErrInvalidParent
	}

	return &Folder{
		id: uuid.New(),
		Metadata: Metadata{
			driveId:    driveId,
			parent:     parent,
			name:       name,
			state:      StatePrivate,
			lastChange: time.Now(),
		},
	}, nil
}

func NewFolderFromRepository(
	id uuid.UUID,
	driveId uuid.UUID,
	parent Parent,
	name string,
	state State,
	lastChange time.Time,
) (*Folder, error) {
	return &Folder{
		id: id,
		Metadata: Metadata{
			driveId:    driveId,
			parent:     parent,
			name:       name,
			state:      state,
			lastChange: lastChange,
		},
	}, nil
}

func (folder *Folder) Id() uuid.UUID {
	return folder.id
}
