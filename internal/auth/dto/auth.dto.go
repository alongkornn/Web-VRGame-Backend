package dto

type RegisterDTO struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Class     string `json:"class"`
	Number    string `json:"number"`
}

type LoginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
