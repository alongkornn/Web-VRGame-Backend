package dto

import "github.com/alongkornn/Web-VRGame-Backend/internal/auth/models"

type Approved struct {
	Status models.Status `json:"status" firestore:"status"`
}

type RoleDTO struct {
	UserId string      `json:"userId" firestore:"userId"`
	Role   models.Role `json:"role" firestore:"role"`
}

type UpdateDTO struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Class     string `json:"class"`
	Number    string `json:"number"`
}

type UpdatePasswordDTO struct {
	Password    string `json:"password"`
	NewPassword string `json:"new_password"`
}
