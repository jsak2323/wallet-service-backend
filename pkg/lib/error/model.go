package error

type Error struct {
	Title   string `json:"title"`
	Message string `json:"message"`
	Trace   string `json:"-"`
}

func NewError(title, message, trace string) *Error {
	return &Error{Message: message, Title: title, Trace: trace}
}

var (
	InternalServerErr = NewError("Internal Server Error", "Error", "")

	ErrorUnmarshalBodyRequest = NewError("error unmarshal body request", "Error", "")
	UsernameNotFound          = NewError("username not found", "Error", "")
	Permissions               = NewError("permissions", "Error", "")
	IncorrectPassword         = NewError("incorrect password", "Error", "")
	InvalidCurrency           = NewError("invalid currency", "Error", "")

	// user
	FailedCreateToken      = NewError("failed create token", "Error", "")
	FailedParseJWT         = NewError("failed parse jwt", "Error", "")
	FailedDeactivateUser   = NewError("failed deactivate user", "Error", "")
	FailedGetAllUser       = NewError("failed get all user", "Error", "")
	FailedCreateUser       = NewError("failed create user", "Error", "")
	FailedUpdateUser       = NewError("failed update user", "Error", "")
	FailedCreateRoleUser   = NewError("failed create role user", "Error", "")
	FailedDeleteRoleUser   = NewError("failed delete role user", "Error", "")
	FailedActivateUser     = NewError("failed activate user", "Error", "")
	FailedGeneratePassword = NewError("failed generate password", "Error", "")

	// role
	FailedGetAllRole       = NewError("failed get all role", "Error", "")
	FailedGetRolesByUserID = NewError("failed get roles by user id", "Error", "")
	FailedGetRoleByUserID  = NewError("failed get role by user id", "Error", "")
	FailedCreateRole       = NewError("failed create role", "Error", "")
	FailedUpdateRole       = NewError("failed update role", "Error", "")
	FailedDeleteRole       = NewError("failed delete role", "Error", "")

	// deposit
	FailedParseFilter     = NewError("failed parse filter", "Error", "")
	FailedParsePagination = NewError("failed parse pagination", "Error", "")
	FailedGetDeposit      = NewError("failed get deposit", "Error", "")

	// role permission
	FailedCreateRolePermission               = NewError("failed create permission", "Error", "")
	FailedDeleteRolePermission               = NewError("failed delete role permission", "Error", "")
	FailedDeleteRolePermissionByPermissionID = NewError("failed delete role permission by permission id", "Error", "")

	// permission
	FailedGetAllPermission     = NewError("failed get all permission", "Error", "")
	FailedGetPermissionByRole  = NewError("failed get permission by role", "Error", "")
	FailedCreatePermission     = NewError("failed create permission", "Error", "")
	FailedUpdatePermission     = NewError("failed update permission", "Error", "")
	FailedDeletePermissionByID = NewError("failed delete permission by id", "Error", "")

	// currency config
	FailedActivateCurrencyConfig   = NewError("failed activate currency config", "Error", "")
	FailedDeactivateCurrencyConfig = NewError("failed deactivate currency config", "Error", "")
	FailedGetAllCurrencyConfig     = NewError("failed get all currency config", "Error", "")
	FailedCreateCurrencyConfig     = NewError("failed create currency config", "Error", "")
	FailedUpdateCurrencyConfig     = NewError("failed update currency config", "Error", "")

	// currency rpc
	FailedCreateCurrencyRPC = NewError("failed create currency rpc", "Error", "")
	FailedDeleteCurrencyRPC = NewError("failed delete currency rpc", "Error", "")

	// rpc config
	FailedGetRPCConfigByCurrencyID = NewError("failed get rpc config by currency id", "Error", "")
	FailedGetRPCConfigByID         = NewError("failed get rpc config by id", "Error", "")
	FailedGetAllRPCConfig          = NewError("failed get all rpc config", "Error", "")
	FailedCreateRPCConfig          = NewError("failed create rpc config", "Error", "")
	FailedUpdateRPCConfig          = NewError("failed update rpc config", "Error", "")
	FailedActivateRPCConfig        = NewError("failed activate rpc config", "Error", "")
	FailedDeactivateRPCConfig      = NewError("failed deactivate rpc config", "Error", "")

	// rpc config rpc method
	FailedCreateRPCConfigRPCMethod              = NewError("failed create rpc config rpc method", "Error", "")
	FailedDeleteRPCConfigRPCMethod              = NewError("failed delete rpc config rpc method", "Error", "")
	FailedDeleteRPCConfigRPCMethodByRPCMethodID = NewError("failed delete rpc config rpc method by rpc method id", "Error", "")

	// rpc method
	FailedGetAllRPCMethod        = NewError("failed get all rpc method", "Error", "")
	FailedGetRPCMethodByConfigID = NewError("failed get rpc method by config id", "Error", "")
	FailedCreateRPCMethod        = NewError("failed create rpc method", "Error", "")
	FailedUpdateRPCMethod        = NewError("failed update rpc method", "Error", "")
	FailedDeleteRPCMethodByID    = NewError("failed delete rpc method by id", "Error", "")

	// rpc request
	FailedGetRPCRequestByRPCMethodID = NewError("failed get rpc request by rpc method id", "Error", "")
	FailedCreateRPCRequest           = NewError("failed create rpc request", "Error", "")
	FailedUpdateRPCRequest           = NewError("failed update rpc request", "Error", "")
	FailedDeleteRPCRequestByID       = NewError("failed delete rpc request by id", "Error", "")

	// rpc response
	FailedGetRPCResponseByRPCMethodID = NewError("failed get rpc response by rpc method id", "Error", "")
	FailedCreateRPCResponse           = NewError("failed create rpc response", "Error", "")
	FailedUpdateRPCResponse           = NewError("failed update rpc response", "Error", "")
	FailedDeleteRPCResponseByID       = NewError("failed delete rpc response by id", "Error", "")

	// cron
	FailedCheckUserBalance          = NewError("failed user balance", "Error", "")
	FailedSendUserBalanceAlertEmail = NewError("send user balance alert email", "Error", "")
	FailedSendReportEmail           = NewError("failed send report email", "Error", "")
	FailedSendHotLimitAlertEmail    = NewError("failed send hot limit alert email", "Error", "")
	FailedCheckHotLimit             = NewError("failed check hot limit", "Error", "")

	ErrorHealthCheckHandler = NewError("error health check handler", "Error", "")

	// balance
	FailedSetColdBalanceDetails = NewError("failed set cold balance details", "Error", "")
	FailedSetHotBalanceDetails  = NewError("failed set hot balance details", "Error", "")
	FailedSetUserBalanceDetails = NewError("failed set user balance details", "Error", "")
	FailedSetPendingWithdraw    = NewError("failed set pending withdraw", "Error", "")
	FailedSetPercent            = NewError("failed set percent", "Error", "")
	FailedSetHotLimits          = NewError("failed set hot limits", "Error", "")

	// wallet cold
	FailedDeactivatedColdBalance   = NewError("failed deactivated cold balance", "Error", "")
	FailedActivatedColdBalance     = NewError("failed activated cold balance", "Error", "")
	FailedCoinToRaw                = NewError("failed coin to raw", "Error", "")
	FailedRawToCoin                = NewError("failed raw to coin", "Error", "")
	FailedCreateColdBalance        = NewError("failed create cold balance", "Error", "")
	FailedUpdateColdBalance        = NewError("failed update cold balance", "Error", "")
	FailedCreateTransaction        = NewError("failed create transaction", "Error", "")
	FailedGetCurrencyByID          = NewError("failed get currency by id", "Error", "")
	FailedFireblocksVaultAccountId = NewError("failed fireblocks vault account id", "Error", "")
	FailedGetVaultAccountAsset     = NewError("failed get vault account asset", "Error", "")
	FailedGetAllColdBalance        = NewError("failed get all cold balance", "Error", "")

	// wallet user
	FailedGetUserBalanceRes = NewError("failed get user balance res", "Error", "")

	// address type
	FailedGetCurrencyBySymbolTokenType = NewError("failed get currency by symbol token type", "Error", "")
	FailedGetRpcConfigByType           = NewError("failed get rpc config by type", "Error", "")
	FailedGetModule                    = NewError("failed get module", "Error", "")
	FailedGetBalance                   = NewError("failed get balance", "Error", "")
	FailedAddressType                  = NewError("failed address type", "Error", "")

	FailedGetMaintenanceList = NewError("failed get maintenance list", "Error", "")
	FailedGetBlockCount      = NewError("failed get block count", "Error", "")
	FailedGetByRpcConfigId   = NewError("failed get by rpc config id", "Error", "")

	FailedGetLogFile = NewError("failed get log file", "Error", "")

	FailedGetNewAddress = NewError("failed get new address", "Error", "")
	FailedSendToAddress = NewError("failed send to address", "Error", "")

	FailedListTransactions   = NewError("failed list transactions", "Error", "")
	FailedUpdateSystemConfig = NewError("failed update system config", "Error", "")

	// withdraws
	FailedGetListWithdraws = NewError("failed get list withdraws", "Error", "")

	InvalidToken   = NewError("invalid token", "Error", "")
	InvalidRequest = NewError("invalid request", "Error", "")
)
