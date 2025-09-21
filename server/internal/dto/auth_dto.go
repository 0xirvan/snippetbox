package dto

type RegisterRequest struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8,max=50"`
}

type LoginRequest struct {
	Email    string `validate:"required"`
	Password string `validate:"required"`
}
