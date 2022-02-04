package user

type User struct {
	Id              int         `json:"id"`
	Username        string      `json:"username" validate:"required"`
	Name            string      `json:"name" validate:"required"`
	Email           string      `json:"email"`
	Password        string      `json:"password" validate:"required,min=8"`
	IPAddress       string      `json:"ip_address"`
	Active          bool        `json:"active"`
	Roles           interface{} `json:"roles"`
	RoleNames       []string    `json:"role_names"`
	PermissionNames []string    `json:"permissions"`
}
