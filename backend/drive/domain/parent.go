package domain

import (
	"errors"
	"github.com/google/uuid"
)

var (
	ErrInvalidParent = errors.New("invalid parent")
)

var (
	ParentKindDrive  = "Drive"
	ParentRecycleBin = "RecycleBin"
	ParentKindFolder = "Folder"
)

type Parent struct {
	kind     string
	folderId uuid.UUID
}

func CreateDriveParent() Parent {
	return Parent{
		kind: ParentKindDrive,
	}
}

func CreateRecycleBinParent() Parent {
	return Parent{
		kind: ParentRecycleBin,
	}
}

func CreateFolderParent(folderId uuid.UUID) (Parent, error) {
	if folderId == uuid.Nil {
		return Parent{}, ErrInvalidParent
	}

	return Parent{kind: ParentKindFolder, folderId: folderId}, nil
}

func CreateParent(value string) (Parent, error) {
	if value == ParentKindDrive {
		return Parent{kind: ParentKindDrive}, nil
	}
	if value == ParentRecycleBin {
		return Parent{kind: ParentRecycleBin}, nil
	}

	folderId, err := uuid.Parse(value)
	if err != nil {
		return Parent{}, ErrInvalidParent
	}

	return Parent{kind: ParentKindFolder, folderId: folderId}, nil
}

func (parent Parent) IsZero() bool {
	return parent == Parent{}
}

func (parent Parent) String() string {
	if parent.folderId == uuid.Nil {
		return parent.kind
	}

	return parent.folderId.String()
}

func (parent Parent) IsDrive() bool {
	return parent.kind == ParentKindDrive
}

func (parent Parent) IsRecycleBin() bool {
	return parent.kind == ParentRecycleBin
}

func (parent Parent) IsFolder() bool {
	return parent.kind == ParentKindFolder
}

func (parent Parent) FolderId() uuid.UUID {
	return parent.folderId
}
