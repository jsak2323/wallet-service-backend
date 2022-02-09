package user

import (
	"context"

	userHandler "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/user"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/jwt"
	"golang.org/x/crypto/bcrypt"
)

func (s *userService) Login(ctx context.Context, loginReq userHandler.LoginReq) (res userHandler.LoginRes, err error) {
	var td jwt.TokenDetails

	if err = s.validator.Validate(loginReq); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.InvalidRequest)
		return res, err
	}

	user, err := s.userRepo.GetByUsername(ctx, loginReq.Username)
	if err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.UsernameNotFound)
		return res, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password)); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.IncorrectPassword)
		return res, err
	}

	if user.RoleNames, err = s.roleRepo.GetNamesByUserId(ctx, user.Id); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedGetRolesByUserID)
		return res, err
	}

	if user.PermissionNames, err = s.permissionRepo.GetNamesByUserId(ctx, user.Id); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.Permissions)
		return res, err
	}

	if td, err = jwt.CreateToken(user); err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedCreateToken)
		return res, err
	}

	res.AccessToken = td.AccessToken
	res.RefreshToken = td.RefreshToken
	return res, nil
}
