package types

// RPC

type GetUserRequest struct {
	SessionId string
}

type GetUserResponse struct {
	Account      string
	Username     string
	Introduction string
}

// HTTP

type LoginRequest struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Account  string `json:"account"`
	Password string `json:"password"`
	Username string `json:"username"`
}
