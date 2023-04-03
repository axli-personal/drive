package domain

import (
	"github.com/google/uuid"
	"time"
)

type File struct {
	Id            uuid.UUID
	ChannelId     uuid.UUID
	SenderAccount string
	Name          string
	Hash          string
	Size          int
	SendAt        time.Time
}

func NewChannelFile(
	channelId uuid.UUID,
	senderAccount string,
	name string,
	hash string,
	size int,
) (*File, error) {
	return &File{
		Id:            uuid.New(),
		ChannelId:     channelId,
		SenderAccount: senderAccount,
		Name:          name,
		Hash:          hash,
		Size:          size,
		SendAt:        time.Now(),
	}, nil
}
