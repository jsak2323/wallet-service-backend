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
	RolesNotFound             = NewError("roles not found", "Error", "")
	Permissions               = NewError("permissions", "Error", "")
	IncorrectPassword         = NewError("incorrect password", "Error", "")

	FailedCreateToken      = NewError("failed create token", "Error", "")
	FailedParseJWT         = NewError("failed parse jwt", "Error", "")
	FailedDeactivateUser   = NewError("failed deactivate user", "Error", "")
	FailedCreateUser       = NewError("failed create user", "Error", "")
	FailedUpdateUser       = NewError("failed update user", "Error", "")
	FailedCreateRoleUser   = NewError("failed create role user", "Error", "")
	FailedDeleteRoleUser   = NewError("failed delete role user", "Error", "")
	FailedActivateUser     = NewError("failed activate user", "Error", "")
	FailedGetAllUser       = NewError("failed get all user", "Error", "")
	FailedGetRoleByUserId  = NewError("failed get role by user id", "Error", "")
	FailedGeneratePassword = NewError("failed generate password", "Error", "")

	InvalidToken   = NewError("invalid token", "Error", "")
	InvalidRequest = NewError("invalid request", "Error", "")
)
