package adapters

import (
	"context"
	"github.com/axli-personal/drive/internal/user/domain"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"time"
)

const (
	SessionPrefix = "session:"
)

type SessionRepository struct {
	client *redis.Client
}

func NewSessionRepository(connectionString string) (SessionRepository, error) {
	options, err := redis.ParseURL(connectionString)
	if err != nil {
		return SessionRepository{}, err
	}

	client := redis.NewClient(options)

	err = client.Ping(context.Background()).Err()
	if err != nil {
		return SessionRepository{}, err
	}

	return SessionRepository{client: client}, nil
}

func (repo SessionRepository) SaveSession(ctx context.Context, session *domain.Session, expire time.Duration) error {
	model := NewSessionModel(session)

	err := repo.client.HSet(ctx, SessionPrefix+model.Id, &model).Err()
	if err != nil {
		return err
	}

	err = repo.client.Expire(ctx, SessionPrefix+model.Id, expire).Err()
	if err != nil {
		return err
	}

	return nil
}

func (repo SessionRepository) GetSession(ctx context.Context, id uuid.UUID) (*domain.Session, error) {
	model := SessionModel{Id: id.String()}

	err := repo.client.HGetAll(ctx, SessionPrefix+model.Id).Scan(&model)
	if err != nil {
		return nil, err
	}

	return model.Session()
}
