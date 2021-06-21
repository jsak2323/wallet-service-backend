package permission

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

func (req DeleteReq) valid() bool {
	if req.Id == 0 {
		return false
	}

	return true
}
