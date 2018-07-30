package messages

import "github.com/opencap/opencap/pkg/types"

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type LookupResponse struct {
	SubType    types.SubTypeId        `json:"sub_type"`
	Address    string                 `json:"address"`
	Extensions map[string]interface{} `json:"extensions"`
}

type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Domain   string `json:"domain"`
	Password string `json:"password"`
}

type LoginResponse struct {
	JWT string `json:"jwt"`
}

type UpdateAddressRequest struct {
	SubType types.SubTypeId `json:"sub_type"`
	Address string          `json:"address"`
}

type AssociateDomainRequest struct {
	Domain string `json:"domain"`
}
