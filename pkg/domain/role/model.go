package role

type Role struct {
	Id    		   int		    `json:"id"`
	Name  		   string	    `json:"name"`
	PermissionList []string     `json:"permission_list"`
	Permissions    interface{}  `json:"permissions"`
}
