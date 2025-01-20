package utils

import (
	"cloud.google.com/go/firestore"
	"github.com/alongkornn/Web-VRGame-Backend/config"
	"github.com/alongkornn/Web-VRGame-Backend/internal/auth/models"
)

func HasAdmin(adminId string) firestore.Query {
	hasAdmin := config.DB.Collection("User").
		Where("is_deleted", "==", false).
		Where("role", "==", models.Admin).
		Where("status", "==", "approved").
		Where("id", "==", adminId).
		Limit(1)

	return hasAdmin
}
