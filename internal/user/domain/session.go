package domain

import (
	"errors"
	"github.com/google/uuid"
)

var (
	ErrMissingUsername = errors.New("missing username")
	ErrMissingId       = errors.New("missing id")
)

type Session struct {
	id       uuid.UUID
	account  Account
	username string
}

func NewSession(account Account, username string) (*Session, error) {
	if account.IsZero() {
		return nil, ErrMissingAccount
	}
	if username == "" {
		return nil, ErrMissingUsername
	}

	return &Session{
		id:       uuid.New(),
		account:  account,
		username: username,
	}, nil
}

func NewSessionFromRepository(id uuid.UUID, account Account, username string) (*Session, error) {
	if id == uuid.Nil {
		return nil, ErrMissingId
	}
	if account.IsZero() {
		return nil, ErrMissingAccount
	}
	if username == "" {
		return nil, ErrMissingUsername
	}

	return &Session{
		id:       id,
		account:  account,
		username: username,
	}, nil
}

func (session *Session) Id() uuid.UUID {
	return session.id
}

func (session *Session) Account() Account {
	return session.account
}

func (session *Session) Username() string {
	return session.username
}
