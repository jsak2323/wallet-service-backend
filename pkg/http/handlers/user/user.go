package user

import (
	"github.com/btcid/wallet-services-backend-go/pkg/domain/permission"
	"github.com/btcid/wallet-services-backend-go/pkg/domain/role"
	"github.com/btcid/wallet-services-backend-go/pkg/domain/user"
	urDmn "github.com/btcid/wallet-services-backend-go/pkg/domain/userrole"
)

type UserService struct {
	userRepo     	user.Repository
	roleRepo     	role.Repository
	urRepo 	 	    urDmn.Repository
	permissionRepo  permission.Repository
}

func NewUserService(
		userRepo user.Repository,
		roleRepo role.Repository,
		urRepo urDmn.Repository,
		permissionRepo permission.Repository,
	) *UserService {
	return &UserService{
		userRepo,
		roleRepo,
		urRepo,
		permissionRepo,
	}
}
