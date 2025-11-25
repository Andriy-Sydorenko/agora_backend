package auth

// TODO: find a more structured and shared way for validating fields in DTOs

type RegisterRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type BasicLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
