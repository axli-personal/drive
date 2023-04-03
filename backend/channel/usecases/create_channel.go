package usecases

import (
	"context"
	"github.com/axli-personal/drive/backend/channel/domain"
	"github.com/axli-personal/drive/backend/channel/repository"
	"github.com/axli-personal/drive/backend/common/decorator"
	"github.com/axli-personal/drive/backend/pkg/errors"
	"github.com/sirupsen/logrus"
)

type (
	CreateChannelArgs struct {
		UserAccount string
		ChannelName string
	}

	CreateChannelResult struct {
	}
)

type createChannelHandler struct {
	channelRepo repository.ChannelRepository
}

func (h createChannelHandler) Handle(ctx context.Context, args CreateChannelArgs) (CreateChannelResult, error) {
	channel, err := domain.NewChannel(args.ChannelName, args.UserAccount)
	if err != nil {
		return CreateChannelResult{}, errors.New(ErrCodeInternal, "fail to create channel", err)
	}

	err = h.channelRepo.SaveChannel(ctx, channel)
	if err != nil {
		return CreateChannelResult{}, errors.New(ErrCodeInternal, "fail to create channel", err)
	}

	return CreateChannelResult{}, nil
}

type CreateChannelHandler decorator.Handler[CreateChannelArgs, CreateChannelResult]

func NewCreateChannelHandler(
	channelRepo repository.ChannelRepository,
	logger *logrus.Entry,
) CreateChannelHandler {
	return decorator.WithLogging[CreateChannelArgs, CreateChannelResult](
		createChannelHandler{
			channelRepo: channelRepo,
		},
		logger,
	)
}
