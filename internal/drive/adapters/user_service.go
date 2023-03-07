package adapters

import (
	"context"
	"github.com/axli-personal/drive/internal/drive/domain"
	"github.com/axli-personal/drive/internal/drive/remote"
	"github.com/axli-personal/drive/internal/pkg/types"
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
		return domain.User{}, err
	}

	return domain.NewUserFromService(response.Account), nil
}
