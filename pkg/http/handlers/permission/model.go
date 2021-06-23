package permission

import (
	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/permission"
)

type CreateReq struct {
	Name string `json:"name"`
}

type CreateRes struct {
	Id 		int 	`json:"id"`
	Message string  `json:"message"`
	Error 	string  `json:"error"`
}

type UpdateReq struct {
	domain.Permission
}

type DeleteReq struct {
	domain.Permission
}

type ListRes struct {
	Permissions []domain.Permission `json:"permissions"`
	Error 		string 		 		`json:"error"`
}

type StandardRes struct {
    Success bool	`json:"success"`
	Message string  `json:"message"`
    Error   string 	`json:"error"`
}