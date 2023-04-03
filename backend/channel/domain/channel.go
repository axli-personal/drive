package domain

import (
	"github.com/google/uuid"
	"time"
)

type Channel struct {
	Id           uuid.UUID
	OwnerAccount string
	Name         string
	InviteToken  string
	CreatedAt    time.Time
}

func NewChannel(name string, ownerAccount string) (*Channel, error) {
	return &Channel{
		Id:           uuid.New(),
		OwnerAccount: ownerAccount,
		Name:         name,
		InviteToken:  uuid.New().String(),
		CreatedAt:    time.Now(),
	}, nil
}
