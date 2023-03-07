package domain

import (
	"errors"
)

var (
	ErrMissingAccount  = errors.New("missing account")
	ErrMissingPassword = errors.New("missing password")
)

type User struct {
	account      Account
	password     Password
	username     string
	introduction string
}

func NewUser(account Account, password Password, username string) (*User, error) {
	if account.IsZero() {
		return nil, ErrMissingAccount
	}
	if password.IsZero() {
		return nil, ErrMissingPassword
	}

	return &User{
		account:  account,
		password: password,
		username: username,
	}, nil
}

func NewUserFromRepository(
	account Account,
	password Password,
	username string,
	introduction string,
) (*User, error) {
	return &User{
		account:      account,
		username:     username,
		password:     password,
		introduction: introduction,
	}, nil
}

func (user *User) Account() Account {
	return user.account
}

func (user *User) Password() Password {
	return user.password
}

func (user *User) Username() string {
	return user.username
}

func (user *User) Introduction() string {
	return user.introduction
}
