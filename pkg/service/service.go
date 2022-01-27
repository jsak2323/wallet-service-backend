package service

// type MysqlRepositories struct {
// 	ColdBalance        coldbalance.Repository
// 	CurrencyConfig     currencyconfig.Repository
// 	CurrencyRpc        currencyrpc.Repository
// 	Deposit            deposit.Repository
// 	HealthCheck        healthcheck.Repository
// 	Permission         permission.Repository
// 	Role               role.Repository
// 	RolePermission     rolepermission.Repository
// 	RpcConfig          rpcconfig.Repository
// 	RpcConfigRpcMethod rpcconfigrpcmethod.Repository
// 	RpcMethod          rpcmethod.Repository
// 	RpcRequest         rpcrequest.Repository
// 	RpcResponse        rpcresponse.Repository
// 	SystemConfig       systemconfig.Repository
// 	User               user.Repository
// 	UserBalance        userbalance.Repository
// 	UserRole           userrole.Repository
// 	Withdraw           withdraw.Repository
// }

// func NewMysqlRepositories(localDB *sql.DB, exchangeSlaveDB *sql.DB) MysqlRepositories {
// 	return MysqlRepositories{
// 		ColdBalance:        NewMysqlColdBalanceRepository(localDB),
// 		CurrencyConfig:     NewMysqlCurrencyConfigRepository(localDB),
// 		CurrencyRpc:        NewMysqlCurrencyRpcRepository(localDB),
// 		Deposit:            NewMysqlDepositRepository(localDB),
// 		HealthCheck:        NewMysqlHealthCheckRepository(localDB),
// 		Permission:         NewMysqlPermissionRepository(localDB),
// 		Role:               NewMysqlRoleRepository(localDB),
// 		RolePermission:     NewMysqlRolePermissionRepository(localDB),
// 		RpcConfig:          NewMysqlRpcConfigRepository(localDB),
// 		RpcConfigRpcMethod: NewMysqlRpcConfigRpcMethodRepository(localDB),
// 		RpcMethod:          NewMysqlRpcMethodRepository(localDB),
// 		RpcRequest:         NewMysqlRpcRequestRepository(localDB),
// 		RpcResponse:        NewMysqlRpcResponseRepository(localDB),
// 		SystemConfig:       NewMysqlSystemConfigRepository(localDB),
// 		User:               NewMysqlUserRepository(localDB),
// 		UserBalance:        NewMysqlUserBalanceRepository(exchangeSlaveDB),
// 		UserRole:           NewMysqlUserRoleRepository(localDB),
// 		Withdraw:           NewMysqlWithdrawRepository(exchangeSlaveDB),
// 	}
// }
