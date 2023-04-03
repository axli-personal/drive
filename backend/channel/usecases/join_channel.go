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

const (
	ErrCodeInvalidInviteToken = "InvalidInviteToken"
)

type (
	JoinChannelArgs struct {
		UserAccount string
		ChannelId   uuid.UUID
		InviteToken string
	}

	JoinChannelResult struct {
	}
)

type joinChannelHandler struct {
	channelRepo repository.ChannelRepository
	memberRepo  repository.MemberRepository
}

func (h joinChannelHandler) Handle(ctx context.Context, args JoinChannelArgs) (JoinChannelResult, error) {
	channel, err := h.channelRepo.GetChannel(ctx, args.ChannelId)
	if err != nil {
		return JoinChannelResult{}, errors.New(ErrCodeChannelNotFound, "fail to join channel", err)
	}

	if args.InviteToken != channel.InviteToken {
		return JoinChannelResult{}, errors.New(ErrCodeInvalidInviteToken, "fail to join channel", err)
	}

	member, err := domain.NewChannelMember(args.ChannelId, args.UserAccount, domain.RoleMember)
	if err != nil {
		return JoinChannelResult{}, errors.New(ErrCodeInternal, "fail to join channel", err)
	}

	err = h.memberRepo.SaveMember(ctx, member)
	if err != nil {
		return JoinChannelResult{}, errors.New(ErrCodeInternal, "fail to join channel", err)
	}

	return JoinChannelResult{}, nil
}

type JoinChannelHandler decorator.Handler[JoinChannelArgs, JoinChannelResult]

func NewJoinChannelHandler(
	channelRepo repository.ChannelRepository,
	memberRepo repository.MemberRepository,
	logger *logrus.Entry,
) JoinChannelHandler {
	return decorator.WithLogging[JoinChannelArgs, JoinChannelResult](
		joinChannelHandler{
			channelRepo: channelRepo,
			memberRepo:  memberRepo,
		},
		logger,
	)
}
