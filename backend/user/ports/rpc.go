package ports

import "github.com/axli-personal/drive/backend/user/service"

type RPCServer struct {
	svc service.Service
}

func NewRPCServer(svc service.Service) RPCServer {
	return RPCServer{svc: svc}
}
