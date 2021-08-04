package cold

type SendToHotReq struct {
	FireblocksName  string  `json:"fireblocks_name"`
	Amount  		string  `json:"amount"`
	Memo 			string  `json:"memo"`
}

type UpdateReq struct {
	Id 		int 	`json:"id"`
	Name 	string 	`json:"name"`
	Balance string  `json:"balance"`
}

type StandardRes struct {
    Success bool	`json:"success"`
	Message string  `json:"message"`
    Error   string 	`json:"error"`
}