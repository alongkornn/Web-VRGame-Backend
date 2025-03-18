package config

import (
	"context"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

var DB *firestore.Client

func InitFirebase() {
	credentialsFile := "/Users/alongkorn/Desktop/gamevr-88a69-firebase-adminsdk-ukt0n-7e6e34d649.json"

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	
	opt := option.WithCredentialsFile(credentialsFile)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("Failed to initialize Firebase App: %v", err)
	}

	DB, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
	}

	log.Println("Successfully connected to Firestore!")
}
