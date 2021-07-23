package role

func (req CreateReq) valid() bool {
	if req.Name == "" {
		return false
	}

	return true
}

func (req UpdateReq) valid() bool {
	if req.Id == 0 {
		return false
	}
	
	if req.Name == "" {
		return false
	}

	return true
}

func (req RolePermissionReq) valid() bool {
	if req.RoleId == 0 {
		return false
	}
	
	if req.PermissionId == 0 {
		return false
	}

	return true
}
