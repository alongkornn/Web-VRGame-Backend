package database

import (
	"context"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/alongkornn/Web-VRGame-Backend/internal/checkpoint/models"
	"github.com/google/uuid"
)

func CreateCheckpointIfNotExists(client *firestore.Client, ctx context.Context) {
	// ตรวจสอบว่ามีข้อมูลใน collection "User" แล้วหรือยัง
	checkpointCollection := client.Collection("Checkpoint")
	docs, err := checkpointCollection.Limit(1).Documents(ctx).GetAll()
	if err != nil {
		log.Fatalf("Error checking checkpoint collection: %v", err)
	}

	if len(docs) == 0 {
		newCheckpoint := models.Checkpoints{
			ID:         uuid.New().String(),
			Name:       "ด่านหนึ่ง",
			Category:   models.Projectile,
			MaxScore:   100,
			PassScore:  50,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
			Is_Deleted: false,
		}

		_, _, err := checkpointCollection.Add(ctx, newCheckpoint)
		if err != nil {
			log.Fatalf("Error creating new checkpoint: %v", err)
		}
		log.Println("Checkpoint created successfully")
	} else {
		log.Println("Checkpoint collection already exists")
	}
}
