package utils

import (
	"cloud.google.com/go/firestore"
	"github.com/alongkornn/Web-VRGame-Backend/config"
	"github.com/alongkornn/Web-VRGame-Backend/internal/auth/models"
)

func HasUser(userId string) firestore.Query {
	hasUser := config.DB.Collection("User").
		Where("is_deleted", "==", false).
		Where("role", "==", models.Player).
		Where("status", "==", models.Approved).
		Where("id", "==", userId).
		Limit(1)
	
	return hasUser
}
