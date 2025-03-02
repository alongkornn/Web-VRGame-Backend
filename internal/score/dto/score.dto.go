package dto

type ResponseScore struct {
	Name  string
	Score int
}

type SetScoreDTO struct {
	Score int    `json:"score"`
	Time  string `json:"time"`
}
