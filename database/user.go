package database

import (
	"context"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/alongkornn/Web-VRGame-Backend/internal/auth/models"
	"github.com/google/uuid"
)

func CreateUserIfNotExists(client *firestore.Client, ctx context.Context) {
	// ตรวจสอบว่ามีข้อมูลใน collection "User" แล้วหรือยัง
	userCollection := client.Collection("User")
	docs, err := userCollection.Limit(1).Documents(ctx).GetAll()
	if err != nil {
		log.Fatalf("Error checking user collection: %v", err)
	}

	if len(docs) == 0 {
		newUser := models.User{
			ID:         uuid.New().String(),
			FirstName:  "admin",
			LastName:   "@admin",
			Email:      "admin@axpi.com",
			Password:   "adminaxpi",
			Role:       models.Admin,
			Status:     models.Approved,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
			Is_Deleted: false,
		}

		_, _, err := userCollection.Add(ctx, newUser)
		if err != nil {
			log.Fatalf("Error creating new user: %v", err)
		}
		log.Println("User created successfully")
	} else {
		log.Println("User collection already exists")
	}
}
