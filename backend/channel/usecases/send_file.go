package usecases

import (
	"context"
	"github.com/axli-personal/drive/backend/channel/domain"
	"github.com/axli-personal/drive/backend/channel/repository"
	"github.com/axli-personal/drive/backend/common/decorator"
	"github.com/axli-personal/drive/backend/pkg/errors"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	ErrCodeUnauthorized = "Unauthorized"
)

type (
	SendFileArgs struct {
		UserAccount string
		ChannelId   uuid.UUID
		FileId      string
	}

	SendFileResult struct {
	}
)

type sendFileHandler struct {
	channelRepo repository.ChannelRepository
	memberRepo  repository.MemberRepository
	fileRepo    repository.FileRepository
}

func (h sendFileHandler) Handle(ctx context.Context, args SendFileArgs) (SendFileResult, error) {
	channel, err := h.channelRepo.GetChannel(ctx, args.ChannelId)
	if err != nil {
		return SendFileResult{}, errors.New(ErrCodeChannelNotFound, "fail to send file", err)
	}

	if args.UserAccount != channel.OwnerAccount {
		member, err := h.memberRepo.FindMember(ctx, args.ChannelId, args.UserAccount)
		if err != nil {
			return SendFileResult{}, errors.New(ErrCodeChannelNotFound, "fail to send file", err)
		}

		if !member.CanSendFile() {
			return SendFileResult{}, errors.New(ErrCodeUnauthorized, "fail to send file", err)
		}
	}

	// TODO: name hash size from drive service.

	channelFile := &domain.File{
		Id:            uuid.New(),
		ChannelId:     args.ChannelId,
		SenderAccount: args.UserAccount,
		Name:          "",
		Hash:          "",
		Size:          0,
		SendAt:        time.Now(),
	}

	err = h.fileRepo.SaveFile(ctx, channelFile)
	if err != nil {
		return SendFileResult{}, errors.New(ErrCodeInternal, "fail to send file", err)
	}

	return SendFileResult{}, nil
}

type SendFileHandler decorator.Handler[SendFileArgs, SendFileResult]

func NewSendFileHandler(
	memberRepo repository.MemberRepository,
	fileRepo repository.FileRepository,
	logger *logrus.Entry,
) SendFileHandler {
	return decorator.WithLogging[SendFileArgs, SendFileResult](
		sendFileHandler{
			memberRepo: memberRepo,
			fileRepo:   fileRepo,
		},
		logger,
	)
}
