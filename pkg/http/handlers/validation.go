package handlers

func (req RegisterReq) valid() bool {
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
