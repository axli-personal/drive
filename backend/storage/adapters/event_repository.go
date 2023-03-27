package adapters

import (
	"context"
	"encoding/json"
	"github.com/axli-personal/drive/backend/pkg/events"
	"github.com/axli-personal/drive/backend/storage/repository"
	"github.com/redis/go-redis/v9"
)

type RedisEventRepository struct {
	client *redis.Client
}

func NewRedisEventRepository(connectionString string) (repository.EventRepository, error) {
	options, err := redis.ParseURL(connectionString)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(options)

	err = client.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}

	return RedisEventRepository{client: client}, nil
}

func (repo RedisEventRepository) PublishFileUploaded(ctx context.Context, event events.FileUploaded) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return repo.client.XAdd(
		ctx,
		&redis.XAddArgs{
			Stream: events.StreamFileUploaded,
			Values: map[string]any{
				events.FieldBody: body,
			},
		},
	).Err()
}

func (repo RedisEventRepository) PublishFileDownloaded(ctx context.Context, event events.FileDownloaded) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return repo.client.XAdd(
		ctx,
		&redis.XAddArgs{
			Stream: events.StreamFileDownloaded,
			Values: map[string]any{
				events.FieldBody: body,
			},
		},
	).Err()
}
