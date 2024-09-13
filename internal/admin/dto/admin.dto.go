package dto

type Approved struct {
	Status string `json:"status" firestore:"status"`
}