package role

import (
	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/role"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

type CreateReq struct {
	Name string `json:"name"`
}

type CreateRes struct {
	Id      int         `json:"id"`
	Message string      `json:"message"`
	Error   *errs.Error `json:"error"`
}

type UpdateReq struct {
	domain.Role
}

type ListRes struct {
	Roles []domain.Role `json:"roles"`
	Error *errs.Error   `json:"error"`
}

type RolePermissionReq struct {
	RoleId       int `json:"role_id"`
	PermissionId int `json:"permission_id"`
}

type StandardRes struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Error   *errs.Error `json:"error"`
}
