package utils

import (
	"cloud.google.com/go/firestore"
	"github.com/alongkornn/Web-VRGame-Backend/config"
)

func GetCheckpointID(name string) firestore.Query {
	query := config.DB.Collection("Checkpoint").
		Where("name", "==", name).
		Limit(1)

	return query
}

func GetCheckpointByID(id string) firestore.Query {
	hasCheckpoint := config.DB.Collection("Checkpoint").
		Where("id", "==", id).
		Limit(1)

	return hasCheckpoint
}
