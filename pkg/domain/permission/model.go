package permission

type Permission struct {
	Id   int    `json:"id" validate:"required"`
	Name string `json:"name"`
}
