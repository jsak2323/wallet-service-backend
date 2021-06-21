package user

func (req CreateReq) valid() bool {
	if req.Username == "" {
		return false
	}

	if req.Name == "" {
		return false
	}

	if req.Password == "" || len(req.Password) < 8 {
		return false
	}

	return true
}

func (req UpdateReq) valid() bool {
	if req.Id == 0 {
		return false
	}
	
	if req.Username == "" {
		return false
	}

	if req.Name == "" {
		return false
	}

	return true
}

func (req UserRoleReq) valid() bool {
	if req.UserId == 0 {
		return false
	}

	if req.RoleId == 0 {
		return false
	}

	return true
}
