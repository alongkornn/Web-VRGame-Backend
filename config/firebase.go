package config

import (
    "context"
    "firebase.google.com/go"
    "firebase.google.com/go/db"
    "log"
    _ "golang.org/x/oauth2/google"
    "google.golang.org/api/option"
)

var Client *db.Client

func InitFirebase() {
    // โหลด serviceAccountKey.json
    sa := option.WithCredentialsFile("/Users/alongkorn/Desktop/gamevr-88a69-firebase-adminsdk-ukt0n-90b587a9b3.json")

    // สร้าง App Firebase
    app, err := firebase.NewApp(context.Background(), nil, sa)
    if err != nil {
        log.Fatalf("error initializing app: %v\n", err)
    }

    // เริ่มต้นฐานข้อมูล
    Client, err = app.Database(context.Background())
    if err != nil {
        log.Fatalf("error initializing database client: %v\n", err)
    }
}
