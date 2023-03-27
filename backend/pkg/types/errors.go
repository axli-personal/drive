package types

import (
	"encoding/json"
)

var (
	ErrCodeInternal        = "Internal"
	ErrCodeUnauthenticated = "Unauthenticated"
	ErrCodeUnauthorized    = "Unauthorized"
)

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail"`
}

func (response ErrorResponse) Error() string {
	data, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}

	return string(data)
}

func NewErrorResponseFromRPC(err error) (ErrorResponse, error) {
	response := ErrorResponse{}

	marshalErr := json.Unmarshal([]byte(err.Error()), &response)
	if marshalErr != nil {
		return ErrorResponse{}, marshalErr
	}

	return response, nil
}
