package role

import (
	"context"

	"github.com/btcid/wallet-services-backend-go/pkg/database/mysql"
	"github.com/btcid/wallet-services-backend-go/pkg/domain/permission"
	"github.com/btcid/wallet-services-backend-go/pkg/domain/role"
	rpDomain "github.com/btcid/wallet-services-backend-go/pkg/domain/rolepermission"
	urDomain "github.com/btcid/wallet-services-backend-go/pkg/domain/userrole"
	roleHandler "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/role"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
)

type RoleService interface {
	GetListRoleDetails(ctx context.Context, page int, limit int) (res roleHandler.ListRes, err error)
	CreateRole(ctx context.Context, createReq roleHandler.CreateReq) (id int, err error)
	CreateRolePermission(ctx context.Context, req roleHandler.RolePermissionReq) (err error)
	UpdateRole(ctx context.Context, updateReq roleHandler.UpdateReq) (err error)
	DeleteRole(ctx context.Context, roleId int) (err error)
	DeleteRolePermission(ctx context.Context, roleId int, permissionId int) (err error)
}

type roleService struct {
	validator      util.CustomValidator
	roleRepo       role.Repository
	permissionRepo permission.Repository
	rpRepo         rpDomain.Repository
	urRepo         urDomain.Repository
}

func NewRoleService(validator util.CustomValidator, mysqlRepos mysql.MysqlRepositories) *roleService {
	return &roleService{
		validator,
		mysqlRepos.Role,
		mysqlRepos.Permission,
		mysqlRepos.RolePermission,
		mysqlRepos.UserRole,
	}
}
