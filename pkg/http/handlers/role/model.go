package role

import (
	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/role"
)

type CreateReq struct {
	Name string `json:"name"`
}

type CreateRes struct {
	Id  	int    `json:"id"`
	Error 	string `json:"error"`
}

type UpdateReq struct {
	domain.Role
}

type ListRes struct {
	Roles []domain.Role `json:"roles"`
	Error string 		`json:"error"`
}

type RolePermissionReq struct {
	RoleId       int  `json:"role_id"`
	PermissionId int `json:"permission_id"`
}

type StandardRes struct {
    Success bool	`json:"success"`
    Error   string 	`json:"error"`
}