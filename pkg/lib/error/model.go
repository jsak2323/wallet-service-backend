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

	InvalidToken   = NewError("invalid token", "Error", "")
	InvalidRequest = NewError("invalid request", "Error", "")
)
