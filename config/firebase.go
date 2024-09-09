package config

import (
	"context"
	"log"
	"time"
	"cloud.google.com/go/firestore"
	_ "golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

var DB *firestore.Client

func InitFirebase() {
	var err error
	// โหลด serviceAccountKey.json 
	sa := option.WithCredentialsFile("/Users/alongkorn/Desktop/gamevr-88a69-firebase-adminsdk-ukt0n-a862e722f6.json")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

	DB, err = firestore.NewClient(ctx, "gamevr-88a69", sa)
    if err != nil {
        log.Fatalf("Failed to create Firestore client: %v", err)
    }

	log.Println("Successfully connectd to firestore")
}

