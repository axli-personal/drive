package ports

import (
	"github.com/axli-personal/drive/internal/user/service"
)

type HTTPServer struct {
	svc service.Service
}

func NewHTTPServer(svc service.Service) HTTPServer {
	return HTTPServer{svc: svc}
}
