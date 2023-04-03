package repository

import (
	"context"
	"github.com/axli-personal/drive/backend/channel/domain"
	"github.com/google/uuid"
)

type Page struct {
	From  int
	Limit int
}

type ChannelRepository interface {
	SaveChannel(ctx context.Context, channel *domain.Channel) error
	GetChannel(ctx context.Context, channelId uuid.UUID) (*domain.Channel, error)
	FindChannel(ctx context.Context, account string) ([]*domain.Channel, error)
}

type MemberRepository interface {
	SaveMember(ctx context.Context, user *domain.Member) error
	FindMember(ctx context.Context, channelId uuid.UUID, account string) (*domain.Member, error)
}

type MessageRepository interface {
	SaveMessage(ctx context.Context, message *domain.Message) error
	FindMessage(ctx context.Context, channelId uuid.UUID, page Page) ([]*domain.Message, error)
}

type FileRepository interface {
	SaveFile(ctx context.Context, file *domain.File) error
	FindFile(ctx context.Context, channelId uuid.UUID, page Page) ([]*domain.File, error)
}
