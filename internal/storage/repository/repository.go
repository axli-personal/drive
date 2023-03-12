package repository

import (
	"context"
	"github.com/axli-personal/drive/internal/pkg/events"
	"github.com/axli-personal/drive/internal/storage/domain"
)

type (
	ObjectRepository interface {
		SaveObject(ctx context.Context, object *domain.Object) error
		GetObject(ctx context.Context, id string) (*domain.Object, error)
	}

	CapacityRepository interface {
		DecreaseRequestCapacity(ctx context.Context) error
	}

	EventRepository interface {
		PublishFileUploaded(ctx context.Context, event events.FileUploaded) error
		PublishFileDownloaded(ctx context.Context, event events.FileDownloaded) error
	}
)
