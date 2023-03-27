package domain

import (
	"errors"
	"github.com/google/uuid"
)

const (
	KB = 1024
	MB = 1024 * KB
	GB = 1024 * MB
)

var (
	ErrExceedStorageUsage = errors.New("exceed storage usage")
)

var (
	StoragePlanFree = StoragePlan{name: "Free", maxBytes: 1 * GB}
)

// How many bytes in the drive.
type StorageUsage struct {
	bytes int
}

func NewStorageUsage(bytes int) (StorageUsage, error) {
	return StorageUsage{bytes: bytes}, nil
}

func (usage StorageUsage) Bytes() int {
	return usage.bytes
}

// The storage plan of the drive, used to control usage.
type StoragePlan struct {
	name     string
	maxBytes int
}

func NewStoragePlan(name string, maxBytes int) (StoragePlan, error) {
	return StoragePlan{name: name, maxBytes: maxBytes}, nil
}

func (plan StoragePlan) Name() string {
	return plan.name
}

func (plan StoragePlan) MaxBytes() int {
	return plan.maxBytes
}

type Drive struct {
	id    uuid.UUID
	owner string
	usage StorageUsage
	plan  StoragePlan
}

func NewDrive(owner string) (*Drive, error) {
	return &Drive{
		id:    uuid.New(),
		owner: owner,
		usage: StorageUsage{
			bytes: 0,
		},
		plan: StoragePlanFree,
	}, nil
}

func NewDriveFromRepository(
	id uuid.UUID,
	owner string,
	usage StorageUsage,
	plan StoragePlan,
) (*Drive, error) {
	return &Drive{
		id:    id,
		owner: owner,
		usage: usage,
		plan:  plan,
	}, nil
}

func (drive *Drive) Id() uuid.UUID {
	return drive.id
}

func (drive *Drive) Owner() string {
	return drive.owner
}

func (drive *Drive) Usage() StorageUsage {
	return drive.usage
}

func (drive *Drive) Plan() StoragePlan {
	return drive.plan
}

func (drive *Drive) IncreaseUsage(bytes int) error {
	if drive.usage.bytes+bytes > drive.plan.maxBytes {
		return ErrExceedStorageUsage
	}

	drive.usage.bytes += bytes

	return nil
}
