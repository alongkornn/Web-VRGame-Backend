package config

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"
	_ "golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

var Client *db.Client

func InitFirebase() {
	// โหลด serviceAccountKey.json 
	sa := option.WithCredentialsFile("/Users/alongkorn/Desktop/gamevr-88a69-firebase-adminsdk-ukt0n-a862e722f6.json")

	// สร้าง Firebase App
	app, err := firebase.NewApp(context.Background(), nil, sa)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	// สร้าง Firestore Client
	client, err := app.Firestore(context.Background())
	if err != nil {
		log.Fatalf("error initializing database client: %v\n", err)
	}
	defer client.Close()

	// ทดสอบการเชื่อมต่อกับ Firestore
	_, _, err = client.Collection("testCollection").Add(context.Background(), map[string]interface{}{
		"testField": "testValue",
	})
	if err != nil {
		log.Fatalf("Failed to add data to Firestore: %v\n", err)
	}

	log.Println("Successfully added data to Firestore")
}

