package config

import (
	"context"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/alongkornn/Web-VRGame-Backend/database"
	_ "golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

var DB *firestore.Client

func InitFirebase() {
	var err error
	// โหลด serviceAccountKey.json
	sa := option.WithCredentialsFile("/Users/VR_1/Desktop/gamevr-88a69-firebase-adminsdk-ukt0n-8b4fa2e924.json")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	DB, err = firestore.NewClient(ctx, "gamevr-88a69", sa)
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
	}
	database.CreateUserIfNotExists(DB, ctx)
	database.CreateCheckpointIfNotExists(DB, ctx)

	log.Println("Successfully connectd to firestore")
}
