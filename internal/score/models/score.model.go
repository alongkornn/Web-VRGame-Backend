package models

import "github.com/alongkornn/Web-VRGame-Backend/internal/checkpoint/models"

type Score struct {
	CheckpointName string
	Category models.Category
	Name string
	Score int
}