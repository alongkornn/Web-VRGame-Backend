package utils

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/alongkornn/Web-VRGame-Backend/config"
)

func GetCheckpointID(name string) (string, error) {
	// ทำการ query ข้อมูล
	query := config.DB.Collection("Checkpoint").
		Where("is_deleted", "==", false).
		Where("name", "==", name).
		Limit(1)

	// Execute query และดึงผลลัพธ์
	docs, err := query.Documents(context.Background()).GetAll()
	if err != nil {
		return "", err
	}

	// ตรวจสอบว่ามีเอกสารหรือไม่
	if len(docs) > 0 {
		// ดึง ID ของเอกสารที่ตรงกับ query
		return docs[0].Ref.ID, nil
	}

	return "", nil // ไม่พบเอกสาร
}

func GetCheckpointByID(id string) firestore.Query {
	hasCheckpoint := config.DB.Collection("Checkpoint").
		Where("is_deleted", "==", false).
		Where("id", "==", id).
		Limit(1)

	return hasCheckpoint
}
