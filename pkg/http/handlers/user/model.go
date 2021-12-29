package user

import (
	userDmn "github.com/btcid/wallet-services-backend-go/pkg/domain/user"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

type CreateReq struct {
	userDmn.User
}

type CreateRes struct {
	Id      int    `json:"id"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

type UpdateReq struct {
	userDmn.User
}

type UserRoleReq struct {
	UserId int `json:"user_id"`
	RoleId int `json:"role_id"`
}

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRes struct {
	AccessToken  string     `json:"access_token"`
	RefreshToken string     `json:"refresh_token"`
	Error        errs.Error `json:"error"`
}

type ListRes struct {
	Users []userDmn.User `json:"users"`
	Error string         `json:"error"`
}

type StandardRes struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error"`
}
