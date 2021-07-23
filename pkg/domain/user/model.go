package user

type User struct {
	Id          	int	 	 	`json:"id"`
	Username    	string 	 	`json:"username"`
	Name        	string 	 	`json:"name"`
	Email 			string 		`json:"email"`
	Password    	string 	 	`json:"password"`
	IPAddress   	string 	 	`json:"ip_address"`
	Active 			bool		`json:"active"`
	Roles       	interface{} `json:"roles"`
	RoleNames 	    []string 	`json:"role_names"`
	PermissionNames []string 	`json:"permissions"`
}
