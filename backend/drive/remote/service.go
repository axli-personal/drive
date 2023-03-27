package remote

import (
	"context"
	"github.com/axli-personal/drive/backend/drive/domain"
)

type (
	UserService interface {
		GetUser(ctx context.Context, sessionId string) (domain.User, error)
	}

	StorageCluster interface {
		ChooseStorageEndPoint(ctx context.Context) (string, error)
	}
)
