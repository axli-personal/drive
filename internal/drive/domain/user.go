package domain

type User struct {
	account string
}

func NewUserFromService(account string) User {
	return User{account: account}
}

func (user User) Account() string {
	return user.account
}
