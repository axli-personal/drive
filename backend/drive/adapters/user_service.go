package adapters

import (
	"context"
	stderr "errors"
	"github.com/axli-personal/drive/backend/drive/domain"
	"github.com/axli-personal/drive/backend/drive/remote"
	"github.com/axli-personal/drive/backend/pkg/errors"
	"github.com/axli-personal/drive/backend/pkg/types"
	"net/rpc"
)

type RPCUserService struct {
	client *rpc.Client
}

func NewRPCUserService(address string) (remote.UserService, error) {
	client, err := rpc.DialHTTP("tcp", address)
	if err != nil {
		return nil, err
	}

	return RPCUserService{client: client}, nil
}

func (service RPCUserService) GetUser(ctx context.Context, sessionId string) (domain.User, error) {
	request := types.GetUserRequest{SessionId: sessionId}
	response := types.GetUserResponse{}

	err := service.client.Call("RPCServer.GetUser", &request, &response)
	if err != nil {
		if errResponse, err := types.NewErrorResponseFromRPC(err); err == nil {
			return domain.User{}, errors.New(errResponse.Code, errResponse.Message, stderr.New(errResponse.Detail))
		}
		return domain.User{}, err
	}

	return domain.NewUserFromService(response.Account), nil
}
