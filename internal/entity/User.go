package entity

type UserResponse struct {
	UserID   string `json:"user_id"`
	GoogleID string `json:"-"`
	Username string `json:"username"`
	Password string `json:"-"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
	Bio      string `json:"bio"`
	Github   string `json:"github"`
	Linkedin string `json:"linkedin"`
	Website  string `json:"website"`
}