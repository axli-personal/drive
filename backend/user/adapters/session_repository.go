package adapters

import (
	"context"
	"github.com/axli-personal/drive/backend/pkg/errors"
	"github.com/axli-personal/drive/backend/user/domain"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	SessionPrefix = "session:"
)

type SessionRepository struct {
	logger *logrus.Entry
	client *redis.Client
}

func NewSessionRepository(connectionString string, logger *logrus.Entry) (SessionRepository, error) {
	options, err := redis.ParseURL(connectionString)
	if err != nil {
		return SessionRepository{}, err
	}

	client := redis.NewClient(options)

	err = client.Ping(context.Background()).Err()
	if err != nil {
		return SessionRepository{}, err
	}

	return SessionRepository{
		logger: logger,
		client: client,
	}, nil
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

	result := repo.client.HGetAll(ctx, SessionPrefix+model.Id)
	if result.Err() != nil {
		return nil, errors.New(domain.ErrCodeRepository, "fail to get session", result.Err())
	}

	if len(result.Val()) == 0 {
		return nil, errors.New(domain.ErrCodeNotFound, "session not found", nil)
	}

	err := result.Scan(&model)
	if err != nil {
		return nil, errors.New(domain.ErrCodeUnmarshal, "fail to scan session", err)
	}

	return model.Session()
}
