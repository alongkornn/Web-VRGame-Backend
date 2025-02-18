package config

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/alongkornn/Web-VRGame-Backend/database"
	"github.com/alongkornn/Web-VRGame-Backend/internal/auth/models"
	websocket_services "github.com/alongkornn/Web-VRGame-Backend/internal/websocket/services"
	"github.com/gorilla/websocket"
	_ "golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

var DB *firestore.Client

func InitFirebase() {
	var err error
	// โหลด serviceAccountKey.json
	sa := option.WithCredentialsFile("C:/Users/VR_1/Desktop/alongkorn/gamevr-88a69-firebase-adminsdk-ukt0n-8b4fa2e924.json")

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

// ListenForUserScoreUpdates เฝ้าดูการเปลี่ยนแปลงคะแนนผู้ใช้งาน
func ListenForUserScoreUpdate() {
	ctx := context.Background()
	query := DB.Collection("User").
		Where("is_deleted", "==", false)

	snapshotIterator := query.Snapshots(ctx)
	defer snapshotIterator.Stop()

	for {
		snapshot, err := snapshotIterator.Next()
		if err != nil {
			log.Println("Error listening to Firestore changes:", err)
			continue
		}

		for _, change := range snapshot.Changes {
			if change.Kind == firestore.DocumentModified {
				updatedUser := &models.User{}
				if err := change.Doc.DataTo(updatedUser); err != nil {
					log.Println("Failed to parse document:", err)
					continue
				}
				BroadcastToClients(updatedUser)
			}
		}
	}
}

// broadcastToClients ส่งข้อมูลที่อัปเดตไปยัง WebSocket Clients
func BroadcastToClients(user *models.User) {
	websocket_services.Mutex.Lock()
	defer websocket_services.Mutex.Unlock()

	data, err := json.Marshal(user)
	if err != nil {
		log.Println("Failed to serialize user data:", err)
		return
	}

	for client := range websocket_services.Clients {
		if err := client.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Println("Failed to send message to client:", err)
			client.Close()
			delete(websocket_services.Clients, client)
		}
	}
}
