package permission

import (
	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/permission"
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
	domain.Permission
}

type DeleteReq struct {
	domain.Permission
}

type ListRes struct {
	Permissions []domain.Permission `json:"permissions"`
	Error       *errs.Error         `json:"error"`
}

type StandardRes struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Error   *errs.Error `json:"error"`
}
