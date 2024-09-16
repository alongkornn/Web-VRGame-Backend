package dto

type Approved struct {
	Status string `json:"status" firestore:"status"`
}

type RoleDTO struct {
	Role string `json:"role" firestore:"role"`
}

type UpdateDTO struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Class     string `json:"class"`
	Number    string `json:"number"`
}

type UpdatePasswordDTO struct {
	Password string `json:"password"`
	NewPassword string `json:"new_password"`
}