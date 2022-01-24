package permission

import (
	"github.com/btcid/wallet-services-backend-go/pkg/domain/permission"
	rpDmn "github.com/btcid/wallet-services-backend-go/pkg/domain/rolepermission"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
)

type PermissionService struct {
	permissionRepo permission.Repository
	rpRepo         rpDmn.Repository
	validator      util.CustomValidator
}

func NewPermissionService(permissionRepo permission.Repository, rpRepo rpDmn.Repository, validator util.CustomValidator) *PermissionService {
	return &PermissionService{
		permissionRepo,
		rpRepo,
		validator,
	}
}
