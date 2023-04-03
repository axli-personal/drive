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
	SendMessageArgs struct {
		UserAccount string
		ChannelId   uuid.UUID
		Text        string
	}

	SendMessageResult struct {
	}
)

type sendMessageHandler struct {
	channelRepo repository.ChannelRepository
	memberRepo  repository.MemberRepository
	messageRepo repository.MessageRepository
}

func (h sendMessageHandler) Handle(ctx context.Context, args SendMessageArgs) (SendMessageResult, error) {
	channel, err := h.channelRepo.GetChannel(ctx, args.ChannelId)
	if err != nil {
		return SendMessageResult{}, errors.New(ErrCodeChannelNotFound, "fail to send file", err)
	}

	if args.UserAccount != channel.OwnerAccount {
		member, err := h.memberRepo.FindMember(ctx, args.ChannelId, args.UserAccount)
		if err != nil {
			return SendMessageResult{}, errors.New(ErrCodeChannelNotFound, "fail to send message", err)
		}

		if !member.CanSendMessage() {
			return SendMessageResult{}, errors.New(ErrCodeUnauthorized, "fail to send message", err)
		}
	}

	channelMessage, err := domain.NewChannelMessage(args.ChannelId, args.UserAccount, args.Text)
	if err != nil {
		return SendMessageResult{}, errors.New(ErrCodeInternal, "fail to send message", err)
	}

	err = h.messageRepo.SaveMessage(ctx, channelMessage)
	if err != nil {
		return SendMessageResult{}, errors.New(ErrCodeInternal, "fail to send message", err)
	}

	return SendMessageResult{}, nil
}

type SendMessageHandler decorator.Handler[SendMessageArgs, SendMessageResult]

func NewSendMessageHandler(
	channelRepo repository.ChannelRepository,
	memberRepo repository.MemberRepository,
	messageRepo repository.MessageRepository,
	logger *logrus.Entry,
) SendMessageHandler {
	return decorator.WithLogging[SendMessageArgs, SendMessageResult](
		sendMessageHandler{
			channelRepo: channelRepo,
			memberRepo:  memberRepo,
			messageRepo: messageRepo,
		},
		logger,
	)
}
