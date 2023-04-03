package usecases

import (
	"context"
	"github.com/axli-personal/drive/backend/channel/domain"
	"github.com/axli-personal/drive/backend/channel/repository"
	"github.com/axli-personal/drive/backend/common/decorator"
	"github.com/axli-personal/drive/backend/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	ErrCodeInternal        = "Internal"
	ErrCodeChannelNotFound = "ChannelNotFound"
)

type (
	ListChannelArgs struct {
		UserAccount string
	}

	ListChannelResult struct {
		Channels []*domain.Channel
	}
)

type listChannelHandler struct {
	channelRepo repository.ChannelRepository
}

func (h listChannelHandler) Handle(ctx context.Context, args ListChannelArgs) (ListChannelResult, error) {
	channels, err := h.channelRepo.FindChannel(ctx, args.UserAccount)
	if err != nil {
		return ListChannelResult{}, errors.New(ErrCodeInternal, "fail to find channel", err)
	}

	return ListChannelResult{Channels: channels}, nil
}

type ListChannelHandler decorator.Handler[ListChannelArgs, ListChannelResult]

func NewListChannelHandler(
	channelRepo repository.ChannelRepository,
	logger *logrus.Entry,
) ListChannelHandler {
	return decorator.WithLogging[ListChannelArgs, ListChannelResult](
		listChannelHandler{
			channelRepo: channelRepo,
		},
		logger,
	)
}
