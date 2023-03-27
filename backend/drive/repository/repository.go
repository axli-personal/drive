package repository

import (
	"context"
	"github.com/axli-personal/drive/backend/drive/domain"
	"github.com/axli-personal/drive/backend/pkg/events"
	"github.com/google/uuid"
)

var (
	ErrCodeRepository = "Repository"
	ErrCodeNotFound   = "NotFound"
)

type (
	DriveRepository interface {
		CreateDrive(ctx context.Context, drive *domain.Drive, options CreateDriveOptions) error
		GetDrive(ctx context.Context, id uuid.UUID) (*domain.Drive, error)
		GetDriveByOwner(ctx context.Context, owner string) (*domain.Drive, error)
		UpdateDrive(ctx context.Context, drive *domain.Drive) error
	}

	CreateDriveOptions struct {
		OnlyOneDrive bool
	}
)

type (
	FolderRepository interface {
		SaveFolder(ctx context.Context, folder *domain.Folder) error
		GetFolder(ctx context.Context, id uuid.UUID) (*domain.Folder, error)
		FindFolder(ctx context.Context, options FindFolderOptions) ([]*domain.Folder, error)
		UpdateFolder(ctx context.Context, folder *domain.Folder, options UpdateFolderOptions) error
	}

	FindFolderOptions struct {
		Parent domain.Parent
	}

	UpdateFolderOptions struct {
		MustInSameState     bool         // Deprecated
		MustInState         domain.State // Deprecated
		UpdateChildrenState bool
	}
)

type (
	FileRepository interface {
		SaveFile(ctx context.Context, file *domain.File) error
		GetFile(ctx context.Context, id uuid.UUID) (*domain.File, error)
		FindFile(ctx context.Context, options FindFileOptions) ([]*domain.File, error)
		UpdateFile(ctx context.Context, file *domain.File, options UpdateFileOptions) error
		DeleteFile(ctx context.Context, file *domain.File) error
	}

	FindFileOptions struct {
		Parent domain.Parent
	}

	UpdateFileOptions struct {
		MustInSameState      bool
		MustInState          domain.State
		IncreaseStorageUsage bool
	}
)

type (
	EventRepository interface {
		PublishFolderRemoved(ctx context.Context, event events.FolderRemoved) error
		PublishFileDeleted(ctx context.Context, event events.FileDeleted) error
		GetFileUploaded(ctx context.Context) (events.FileUploaded, error)
		GetFileDownloaded(ctx context.Context) (events.FileDownloaded, error)
		AckGetFileUploaded(ctx context.Context, id string) error
		AckFileDownloaded(ctx context.Context, id string) error
	}
)
