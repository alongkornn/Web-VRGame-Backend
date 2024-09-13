package dto

type Approved struct {
	Status string `json:"status" firestore:"status"`
}

type RoleDTO struct {
	Role string `json:"role" firestore:"role"`
}