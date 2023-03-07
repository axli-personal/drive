package ports

import (
	"context"
	"github.com/axli-personal/drive/internal/pkg/types"
	"github.com/axli-personal/drive/internal/user/usecases"
	"github.com/google/uuid"
)

func (server RPCServer) GetUser(request *types.GetUserRequest, response *types.GetUserResponse) (err error) {
	sessionId, err := uuid.Parse(request.SessionId)
	if err != nil {
		return err
	}

	result, err := server.svc.GetUser.Handle(
		context.Background(),
		usecases.GetUserArgs{
			SessionId: sessionId,
		},
	)
	if err != nil {
		return err
	}

	response.Account = result.Account.String()
	response.Username = result.Username
	response.Introduction = result.Introduction

	return nil
}
