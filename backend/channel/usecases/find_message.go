package usecases

import (
	"context"
	"github.com/axli-personal/drive/backend/channel/domain"
	"github.com/axli-personal/drive/backend/channel/repository"
	"github.com/axli-personal/drive/backend/common/decorator"
	"github.com/axli-personal/drive/backend/pkg/errors"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type (
	FindMessageArgs struct {
		UserAccount string
		ChannelId   uuid.UUID
		From        int
		Limit       int
	}

	FindMessageResult struct {
		Messages []*domain.Message
	}
)

type findMessageHandler struct {
	channelRepo repository.ChannelRepository
	memberRepo  repository.MemberRepository
	messageRepo repository.MessageRepository
}

func (h findMessageHandler) Handle(ctx context.Context, args FindMessageArgs) (FindMessageResult, error) {
	channel, err := h.channelRepo.GetChannel(ctx, args.ChannelId)
	if err != nil {
		return FindMessageResult{}, errors.New(ErrCodeChannelNotFound, "fail to find message", err)
	}

	if args.UserAccount != channel.OwnerAccount {
		member, err := h.memberRepo.FindMember(ctx, args.ChannelId, args.UserAccount)
		if err != nil {
			return FindMessageResult{}, errors.New(ErrCodeChannelNotFound, "fail to find message", err)
		}

		if !member.CanReadMessage() {
			return FindMessageResult{}, errors.New(ErrCodeUnauthorized, "fail to find message", err)
		}
	}

	messages, err := h.messageRepo.FindMessage(
		ctx,
		args.ChannelId,
		repository.Page{
			From:  args.From,
			Limit: args.Limit,
		},
	)
	if err != nil {
		return FindMessageResult{}, errors.New(ErrCodeInternal, "fail to find message", err)
	}

	return FindMessageResult{Messages: messages}, nil
}

type FindMessageHandler decorator.Handler[FindMessageArgs, FindMessageResult]

func NewFindMessageHandler(
	channelRepo repository.ChannelRepository,
	memberRepo repository.MemberRepository,
	messageRepo repository.MessageRepository,
	logger *logrus.Entry,
) FindMessageHandler {
	return decorator.WithLogging[FindMessageArgs, FindMessageResult](
		findMessageHandler{
			channelRepo: channelRepo,
			memberRepo:  memberRepo,
			messageRepo: messageRepo,
		},
		logger,
	)
}
