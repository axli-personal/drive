package ports

import (
	"context"
	"github.com/axli-personal/drive/backend/pkg/errors"
	"github.com/axli-personal/drive/backend/pkg/types"
	"github.com/axli-personal/drive/backend/user/usecases"
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
		if err, ok := err.(*errors.Error); ok {
			if err.Code() == usecases.ErrCodeNotLogin {
				return types.ErrorResponse{
					Code:    types.ErrCodeUnauthenticated,
					Message: "authentication failed",
					Detail:  err.Error(),
				}
			}
		}
		return types.ErrorResponse{
			Code:    types.ErrCodeInternal,
			Message: "fail to get user",
			Detail:  err.Error(),
		}
	}

	response.Account = result.Account.String()
	response.Username = result.Username
	response.Introduction = result.Introduction

	return nil
}
