package userrole

type UserRole struct {
	UserId int
	RoleId int
	Role   interface{}
}

type Relation struct {
	User bool
	Role bool
}