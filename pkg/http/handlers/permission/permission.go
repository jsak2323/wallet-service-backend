package permission

import (
	"github.com/btcid/wallet-services-backend-go/pkg/domain/permission"
	rpDmn "github.com/btcid/wallet-services-backend-go/pkg/domain/rolepermission"
)

type PermissionService struct {
	permissionRepo  permission.Repository
	rpRepo 			rpDmn.Repository
}

func NewPermissionService(permissionRepo permission.Repository, rpRepo rpDmn.Repository) *PermissionService {
	return &PermissionService{
		permissionRepo,
		rpRepo,
	}
}
