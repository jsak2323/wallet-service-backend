package role

type Role struct {
	Id             int         `json:"id" validate:"required"`
	Name           string      `json:"name" validate:"required"`
	PermissionList []string    `json:"permission_list"`
	Permissions    interface{} `json:"permissions"`
}
