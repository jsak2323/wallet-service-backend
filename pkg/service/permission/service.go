package permission

import (
	"context"

	"github.com/btcid/wallet-services-backend-go/pkg/database/mysql"
	"github.com/btcid/wallet-services-backend-go/pkg/domain/permission"
	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/permission"
	rp "github.com/btcid/wallet-services-backend-go/pkg/domain/rolepermission"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
)

type PermissionService interface {
	ListPermissions(ctx context.Context, page int, limit int) (res []permission.Permission, err error)
	CreatePermission(ctx context.Context, name string) (id int, err error)
	UpdatePermission(ctx context.Context, updateReq domain.Permission) (err error)
	DeletePermission(ctx context.Context, permissionId int) (err error)
}

type permissionService struct {
	validator      util.CustomValidator
	permissionRepo permission.Repository
	rpRepo         rp.Repository
}

func NewPermissionService(validator util.CustomValidator, mysqlRepos mysql.MysqlRepositories) *permissionService {
	return &permissionService{
		validator,
		mysqlRepos.Permission,
		mysqlRepos.RolePermission,
	}
}
