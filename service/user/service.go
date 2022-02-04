package user

import (
	"context"

	"github.com/btcid/wallet-services-backend-go/pkg/database/mysql"
	"github.com/btcid/wallet-services-backend-go/pkg/domain/permission"
	"github.com/btcid/wallet-services-backend-go/pkg/domain/role"
	"github.com/btcid/wallet-services-backend-go/pkg/domain/user"
	urDmn "github.com/btcid/wallet-services-backend-go/pkg/domain/userrole"
	domainUser "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/user"
	userHandler "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/user"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
	"github.com/btcid/wallet-services-backend-go/pkg/thirdparty/exchange"
)

type UserService interface {
	Login(ctx context.Context, loginReq userHandler.LoginReq) (res domainUser.LoginRes, err error)
	ActivateUser(ctx context.Context, userId int) (err error)
	DeactivateUser(ctx context.Context, userId int) (err error)
	ListUser(ctx context.Context, page int, limit int) (res []user.User, err error)
	CreateUser(ctx context.Context, createReq domainUser.CreateReq) (id int, err error)
	UpdateUser(ctx context.Context, updateReq userHandler.UpdateReq) (err error)
	CreateUserRole(ctx context.Context, urReq userHandler.UserRoleReq) (err error)
	DeleteUserRole(ctx context.Context, userId int, roleId int) (err error)
}

type userService struct {
	validator      util.CustomValidator
	userRepo       user.Repository
	roleRepo       role.Repository
	urRepo         urDmn.Repository
	permissionRepo permission.Repository
}

func NewUserService(
	validator util.CustomValidator,
	mysqlRepos mysql.MysqlRepositories,
	exchangeApiRepos exchange.APIRepositories,
) *userService {
	return &userService{
		validator,
		mysqlRepos.User,
		mysqlRepos.Role,
		mysqlRepos.UserRole,
		mysqlRepos.Permission,
	}
}
