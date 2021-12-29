package error

type Error struct {
	Title   string `json:"title"`
	Message string `json:"message"`
	Trace   string
}

func NewError(title, message, trace string) Error {
	return Error{Message: message, Title: title, Trace: trace}
}

var (
	InternalServerErr         = NewError("Internal Server Error", "Error", "")
	ErrorUnmarshalBodyRequest = NewError("error unmarshal body request", "Error", "")
	UsernameNotFound          = NewError("username not found", "Error", "")
	RolesNotFound             = NewError("roles not found", "Error", "")
	Permissions               = NewError("permissions", "Error", "")
	IncorrectPassword         = NewError("incorrect password", "Error", "")
	FailedCreateToken         = NewError("failed create token", "Error", "")
)
