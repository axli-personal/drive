package adapters

import (
	"context"
	"encoding/json"
	"github.com/axli-personal/drive/backend/drive/repository"
	"github.com/axli-personal/drive/backend/pkg/events"
	"github.com/redis/go-redis/v9"
)

var (
	GroupDrive = "Drive"
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

	ctx := context.Background()

	err = client.Ping(ctx).Err()
	if err != nil {
		return nil, err
	}

	// Will return BUSYGROUP error when group already exists.
	_ = client.XGroupCreateMkStream(ctx, events.StreamFileUploaded, GroupDrive, "0")
	_ = client.XGroupCreateMkStream(ctx, events.StreamFileDownloaded, GroupDrive, "0")

	return RedisEventRepository{client: client}, nil
}

func (repo RedisEventRepository) PublishFolderRemoved(ctx context.Context, event events.FolderRemoved) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return repo.client.XAdd(
		ctx,
		&redis.XAddArgs{
			Stream: events.StreamFolderRemoved,
			Values: map[string]any{
				events.FieldBody: body,
			},
		},
	).Err()
}

func (repo RedisEventRepository) PublishFileDeleted(ctx context.Context, event events.FileDeleted) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return repo.client.XAdd(
		ctx,
		&redis.XAddArgs{
			Stream: events.StreamFileDeleted,
			Values: map[string]any{
				events.FieldBody: body,
			},
		},
	).Err()
}

func (repo RedisEventRepository) GetFileUploaded(ctx context.Context) (events.FileUploaded, error) {
	result, err := repo.client.XReadGroup(
		ctx,
		&redis.XReadGroupArgs{
			Group:    GroupDrive,
			Consumer: "1",
			Streams:  []string{events.StreamFileUploaded, ">"},
			Count:    1,
			Block:    0,
		},
	).Result()
	if err != nil {
		return events.FileUploaded{}, err
	}

	// TODO: may overflow?
	msg := result[0].Messages[0]

	body, ok := msg.Values[events.FieldBody].(string)
	if !ok {
		return events.FileUploaded{}, err
	}

	event := events.FileUploaded{
		EventId: msg.ID,
	}

	err = json.Unmarshal([]byte(body), &event)
	if err != nil {
		return events.FileUploaded{}, err
	}

	return event, nil
}

func (repo RedisEventRepository) GetFileDownloaded(ctx context.Context) (events.FileDownloaded, error) {
	result, err := repo.client.XReadGroup(
		ctx,
		&redis.XReadGroupArgs{
			Group:    GroupDrive,
			Consumer: "1",
			Streams:  []string{events.StreamFileDownloaded, ">"},
			Count:    1,
			Block:    0,
		},
	).Result()
	if err != nil {
		return events.FileDownloaded{}, err
	}

	// TODO: may overflow?
	msg := result[0].Messages[0]

	body, ok := msg.Values[events.FieldBody].(string)
	if !ok {
		return events.FileDownloaded{}, err
	}

	event := events.FileDownloaded{
		EventId: msg.ID,
	}

	err = json.Unmarshal([]byte(body), &event)
	if err != nil {
		return events.FileDownloaded{}, err
	}

	return event, nil
}

func (repo RedisEventRepository) AckGetFileUploaded(ctx context.Context, id string) error {
	return repo.client.XAck(ctx, events.StreamFileUploaded, GroupDrive, id).Err()
}

func (repo RedisEventRepository) AckFileDownloaded(ctx context.Context, id string) error {
	return repo.client.XAck(ctx, events.StreamFileDownloaded, GroupDrive, id).Err()
}
