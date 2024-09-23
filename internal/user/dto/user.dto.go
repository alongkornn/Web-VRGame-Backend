package dto

type UpdateUserDTO struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Class     string `json:"class"`
	Number    string `json:"number"`
}
