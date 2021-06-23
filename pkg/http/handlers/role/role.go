package role

import (
	"github.com/btcid/wallet-services-backend-go/pkg/domain/permission"
	"github.com/btcid/wallet-services-backend-go/pkg/domain/role"
	rpDomain "github.com/btcid/wallet-services-backend-go/pkg/domain/rolepermission"
	urDomain "github.com/btcid/wallet-services-backend-go/pkg/domain/userrole"
)

const errInternalServer = "Internal server error"

type RoleService struct {
	roleRepo     	role.Repository
	permissionRepo  permission.Repository
	rpRepo			rpDomain.Repository
	urRepo		    urDomain.Repository
}

func NewRoleService(
		roleRepo role.Repository,
		permissionRepo permission.Repository,
		rpRepo rpDomain.Repository,
		urRepo urDomain.Repository,
	) *RoleService {
	return &RoleService{
		roleRepo,
		permissionRepo,
		rpRepo,
		urRepo,
	}
}
