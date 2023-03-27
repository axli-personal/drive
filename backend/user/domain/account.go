package domain

import "errors"

type Account struct {
	value string
}

func NewAccount(account string) (Account, error) {
	if len(account) < 3 || len(account) > 30 {
		return Account{}, errors.New("invalid account length")
	}

	for i := 0; i < len(account); i++ {
		if 'a' <= account[i] && account[i] <= 'z' {
			continue
		}
		if 'A' <= account[i] && account[i] <= 'Z' {
			continue
		}
		if '0' <= account[i] && account[i] <= '9' {
			continue
		}
		if account[i] == '-' || account[i] == '_' {
			continue
		}
		return Account{}, errors.New("invalid account charset")
	}

	return Account{value: account}, nil
}

func (a Account) IsZero() bool {
	return a == Account{}
}

func (a Account) String() string {
	return a.value
}
