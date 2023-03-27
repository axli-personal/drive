package repository

import (
	"context"
	"github.com/axli-personal/drive/backend/pkg/events"
	"github.com/axli-personal/drive/backend/storage/domain"
)

type (
	ObjectRepository interface {
		SaveObject(ctx context.Context, object *domain.Object) error
		GetObject(ctx context.Context, hash string) (*domain.Object, error)
	}

	CapacityRepository interface {
		DecreaseRequestCapacity(ctx context.Context) error
	}

	EventRepository interface {
		PublishFileUploaded(ctx context.Context, event events.FileUploaded) error
		PublishFileDownloaded(ctx context.Context, event events.FileDownloaded) error
	}
)
