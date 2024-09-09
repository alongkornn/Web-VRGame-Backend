package dto

type RegisterDTO struct {
	FirstName string `json:"firstname" validate:"required"`
	LastName  string `json:"lastname" validate:"required"`
	Email     string `json:"email" validate:"required, email"`
	Password  string `json:"password" validate:"required, min=8"`
	Class     string `json:"class" validate:"required"`
	Number    string `json:"number" validate:"required"`
}

type LoginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}