package adapters_test

import (
	"context"
	"github.com/axli-personal/drive/internal/pkg/events"
	"github.com/axli-personal/drive/internal/storage/adapters"
	"github.com/google/uuid"
	"testing"
)

func TestRedisEventRepository(t *testing.T) {
	connectString := "redis://localhost:6379"

	repo, err := adapters.NewRedisEventRepository(connectString)
	if err != nil {
		t.Fatal(err)
	}

	err = repo.PublishFileUploaded(
		context.Background(),
		events.FileUploaded{
			Endpoint:   "https://endpoint",
			FileId:     uuid.New().String(),
			TotalBytes: 100,
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	err = repo.PublishFileDownloaded(
		context.Background(),
		events.FileDownloaded{
			FileId: uuid.New().String(),
		},
	)
	if err != nil {
		t.Fatal(err)
	}
}
