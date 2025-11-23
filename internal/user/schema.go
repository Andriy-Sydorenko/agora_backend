package user

// TODO: find a more structured and shared way for validating fields in DTOs

type PublicUserResponse struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=3,max=50,alphanum"`
}
