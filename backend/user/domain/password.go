package domain

import "errors"

type Password struct {
	value string
}

func NewPassword(password string) (Password, error) {
	if len(password) < 6 || len(password) > 30 {
		return Password{}, errors.New("invalid password length")
	}

	return Password{value: password}, nil
}

func (p Password) IsZero() bool {
	return p == Password{}
}

func (p Password) String() string {
	return p.value
}
