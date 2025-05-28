package user_dto

type UpdateUserRequest struct {
	Username *string `json:"username" validate:"omitempty,min=3,max=30"`
	Password *string `json:"password" validate:"omitempty,min=6,max=12"`
	Email    *string `json:"email" validate:"omitempty,email"`
	Avatar   *string `json:"avatar" validate:"omitempty"`
	Bio      *string `json:"bio" validate:"omitempty,min=5,max=255"`
	Github   *string `json:"github" validate:"omitempty,url"`
	Linkedin *string `json:"linkedin" validate:"omitempty,url"`
	Website  *string `json:"website" validate:"omitempty,url"`
}
