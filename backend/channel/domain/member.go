package domain

import (
	"github.com/google/uuid"
	"time"
)

const (
	RoleAdmin  = "Admin"
	RoleMember = "Member"
)

type Member struct {
	Id        uuid.UUID
	ChannelId uuid.UUID
	Account   string
	Role      string
	JoinAt    time.Time
}

func NewChannelMember(channelId uuid.UUID, account string, role string) (*Member, error) {
	return &Member{
		Id:        uuid.New(),
		ChannelId: channelId,
		Account:   account,
		Role:      role,
		JoinAt:    time.Now(),
	}, nil
}

func (m *Member) CanReadMessage() bool {
	return true
}

func (m *Member) CanSendMessage() bool {
	return true
}

func (m *Member) CanSendFile() bool {
	return m.Role == RoleAdmin
}
