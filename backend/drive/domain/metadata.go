package domain

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

var (
	ErrCannotRead         = errors.New("cannot read")
	ErrCannotReadParent   = errors.New("cannot read parent")
	ErrCannotWrite        = errors.New("cannot write")
	ErrCannotDelete       = errors.New("cannot delete")
	ErrCannotShare        = errors.New("cannot share")
	ErrCannotTrash        = errors.New("cannot trash")
	ErrCannotChangeParent = errors.New("cannot change parent")
)

type Metadata struct {
	driveId    uuid.UUID
	parent     Parent
	name       string
	state      State
	lastChange time.Time
}

func (metadata *Metadata) DriveId() uuid.UUID {
	return metadata.driveId
}

func (metadata *Metadata) Parent() Parent {
	return metadata.parent
}

func (metadata *Metadata) Name() string {
	return metadata.name
}

func (metadata *Metadata) State() State {
	return metadata.state
}

func (metadata *Metadata) LastChange() time.Time {
	return metadata.lastChange
}

func (metadata *Metadata) SetState(state State) {
	metadata.state = state
}

func (metadata *Metadata) CanRead(userDriveId uuid.UUID) error {
	if metadata.state == StateShared {
		return nil
	}
	if metadata.state == StatePrivate && metadata.driveId == userDriveId {
		return nil
	}

	return ErrCannotRead
}

func (metadata *Metadata) CanReadParent(userDriveId uuid.UUID) error {
	if metadata.driveId == userDriveId {
		return nil
	}

	return ErrCannotReadParent
}

func (metadata *Metadata) CanWrite(userDriveId uuid.UUID) error {
	if metadata.driveId != userDriveId {
		return ErrCannotWrite
	}

	return nil
}

func (metadata *Metadata) CanDelete(userDriveId uuid.UUID) error {
	if metadata.state != StateTrashed || metadata.driveId != userDriveId {
		return ErrCannotDelete
	}

	return nil
}

func (metadata *Metadata) Share() error {
	if metadata.state != StatePrivate {
		return ErrCannotShare
	}

	metadata.state = StateShared
	metadata.lastChange = time.Now()

	return nil
}

func (metadata *Metadata) ChangeParent(parent Parent) error {
	if parent.IsZero() || parent.IsRecycleBin() {
		return ErrInvalidParent
	}

	if metadata.state != StatePrivate && metadata.state != StateShared {
		return ErrCannotChangeParent
	}

	if metadata.parent == parent {
		return ErrCannotChangeParent
	}

	metadata.parent = parent
	metadata.lastChange = time.Now()

	return nil
}

func (metadata *Metadata) Trash() error {
	if metadata.state != StatePrivate && metadata.state != StateShared {
		return ErrCannotTrash
	}

	metadata.state = StateTrashed
	metadata.parent = CreateRecycleBinParent()
	metadata.lastChange = time.Now()

	return nil
}
