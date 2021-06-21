package user

type User struct {
	Id          	int	 	 	`json:"id"`
	Username    	string 	 	`json:"username"`
	Name        	string 	 	`json:"name"`
	Password    	string 	 	`json:"password"`
	IPAddress   	string 	 	`json:"ip_address"`
	Roles       	interface{} `json:"roles"`
	RoleNames 	    []string 	`json:"role_names"`
	PermissionNames []string 	`json:"permissions"`
}
