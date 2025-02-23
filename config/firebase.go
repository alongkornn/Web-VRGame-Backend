package config

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

var DB *firestore.Client
var databaseURL = "https://gamevr-88a69-default-rtdb.asia-southeast1.firebasedatabase.app/"

func InitFirebase() {
	credentialsFile := "/Users/alongkorn/Desktop/gamevr-88a69-firebase-adminsdk-ukt0n-7e6e34d649.json"

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	
	opt := option.WithCredentialsFile(credentialsFile)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("❌ Failed to initialize Firebase App: %v", err)
	}

	DB, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalf("❌ Failed to create Firestore client: %v", err)
	}

	log.Println("Successfully connected to Firestore!")
}

func AddToRealtimeDB(userID string, data interface{}) error {
	// ใช้ PUT เพื่ออัปเดตข้อมูลที่ตำแหน่ง users/userID
	url := fmt.Sprintf("%s/users/%s.json", databaseURL, userID)

	// แปลงข้อมูลเป็น JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %v", err)
	}

	// ส่ง HTTP PUT request ไปยัง Realtime Database เพื่ออัปเดตข้อมูลที่ตำแหน่งที่กำหนด
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// ส่งคำขอ HTTP
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// อ่าน Response Body
	body, _ := ioutil.ReadAll(resp.Body)

	// ตรวจสอบว่า Status Code เป็น OK หรือไม่
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error response from Firebase: %s", body)
	}

	log.Println("Successfully added data to Realtime Database!")
	return nil
}

func UpdateStatusInRealtimeDB(id string, status string) error {
	// สร้าง URL ที่จะอัปเดตข้อมูลใน Realtime Database
	url := fmt.Sprintf("%s/users/%s.json", databaseURL, id)

	// ดึงข้อมูลผู้ใช้ปัจจุบันจาก Realtime Database
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to get current data: %v", err)
	}
	defer resp.Body.Close()

	var currentData map[string]interface{}
	body, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &currentData); err != nil {
		return fmt.Errorf("failed to unmarshal current data: %v", err)
	}

	// เพิ่มหรืออัปเดตฟิลด์ status
	currentData["status"] = status

	// แปลงข้อมูลทั้งหมดกลับเป็น JSON
	jsonData, err := json.Marshal(currentData)
	if err != nil {
		return fmt.Errorf("failed to marshal updated data: %v", err)
	}

	// ส่ง HTTP PATCH request ไปยัง Realtime Database
	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// ส่งคำขอ HTTP
	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// อ่าน Response Body
	body, _ = ioutil.ReadAll(resp.Body)

	// ตรวจสอบว่า Status Code เป็น OK หรือไม่
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error response from Firebase: %s", body)
	}

	log.Println("Successfully updated status in Realtime Database!")
	return nil
}

func UpdateCurrentCheckpointInRealtimeDB(userID string, currentCheckpointID string, score int) error {
	// ตรวจสอบว่า currentCheckpointID ไม่เป็นค่าว่าง
	if currentCheckpointID == "" {
		return fmt.Errorf("invalid checkpoint ID")
	}

	// สร้าง URL ที่จะอัปเดตข้อมูลใน Realtime Database
	url := fmt.Sprintf("%s/users/%s.json", databaseURL, userID)

	// สร้างข้อมูลที่ต้องการอัปเดต
	updateData := map[string]interface{}{
		"currentCheckpoint": currentCheckpointID,
		"score":             score,
	}

	// แปลงข้อมูลเป็น JSON
	jsonData, err := json.Marshal(updateData)
	if err != nil {
		return fmt.Errorf("failed to marshal updated data: %v", err)
	}

	// ส่ง HTTP PATCH request ไปยัง Realtime Database (แทน PUT)
	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// ส่งคำขอ HTTP
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// อ่าน Response Body
	body, _ := ioutil.ReadAll(resp.Body)

	// ตรวจสอบว่า Status Code เป็น OK หรือไม่
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error response from Firebase: %s", body)
	}

	log.Println("Successfully updated current_checkpoint in Realtime Database!")
	return nil
}
