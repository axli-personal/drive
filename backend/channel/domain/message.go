package domain

import (
	"github.com/google/uuid"
	"time"
)

type Message struct {
	Id            uuid.UUID
	ChannelId     uuid.UUID
	SenderAccount string
	Text          string
	SendAt        time.Time
}

func NewChannelMessage(channelId uuid.UUID, senderAccount string, text string) (*Message, error) {
	return &Message{
		Id:            uuid.New(),
		ChannelId:     channelId,
		SenderAccount: senderAccount,
		Text:          text,
		SendAt:        time.Now(),
	}, nil
}
