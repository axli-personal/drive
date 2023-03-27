package domain

import (
	"context"
	"github.com/google/uuid"
	"time"
)

var (
	ErrCodeRepository = "Repository"
	ErrCodeNotFound   = "NotFound"
	ErrCodeUnmarshal  = "Unmarshal"
)

type (
	SessionRepository interface {
		SaveSession(ctx context.Context, session *Session, expire time.Duration) error

		GetSession(ctx context.Context, id uuid.UUID) (*Session, error)
	}

	UserRepository interface {
		SaveUser(ctx context.Context, user *User) error

		GetUser(ctx context.Context, account Account) (*User, error)
	}
)
