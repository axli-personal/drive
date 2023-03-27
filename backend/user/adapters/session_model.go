package adapters

import (
	"github.com/axli-personal/drive/backend/user/domain"
	"github.com/google/uuid"
)

type SessionModel struct {
	Id       string `redis:"-"`
	Account  string `redis:"account"`
	Username string `redis:"username"`
}

func NewSessionModel(session *domain.Session) SessionModel {
	return SessionModel{
		Id:       session.Id().String(),
		Account:  session.Account().String(),
		Username: session.Username(),
	}
}

func (model SessionModel) Session() (*domain.Session, error) {
	sessionId, err := uuid.Parse(model.Id)
	if err != nil {
		return nil, err
	}

	account, err := domain.NewAccount(model.Account)
	if err != nil {
		return nil, err
	}

	return domain.NewSessionFromRepository(sessionId, account, model.Username)
}
