package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/alongkornn/Web-VRGame-Backend/config"
	"github.com/alongkornn/Web-VRGame-Backend/internal/auth/dto"
	"github.com/alongkornn/Web-VRGame-Backend/internal/auth/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// type AuthService struct {
// 	db *firestore.Client
// }

// func NewAuthService(db *firestore.Client) *AuthService {
// 	return &AuthService{db: db}
// }

func Register(ctx context.Context, registerDTO *dto.RegisterDTO) (int, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(registerDTO.Password), bcrypt.DefaultCost)
	if err != nil {
		return http.StatusBadRequest, errors.New("hash password is error")
	}
	user := models.User{
		ID:         uuid.New().String(),
		FirstName:  registerDTO.FirstName,
		LastName:   registerDTO.LastName,
		Email:      registerDTO.Email,
		Password:   string(hashPassword),
		Class:      registerDTO.Class,
		Number:     registerDTO.Number,
		Score:      0,
		Level:      1,
		Role:       models.Player,
		Status:     models.Pending,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Is_Deleted: false,
	}

	_, _, err = config.DB.Collection("User").Add(ctx, user)
    if err != nil {
        // Print error to the console for debugging
        fmt.Printf("Error adding document: %v\n", err)
        return http.StatusInternalServerError, errors.New("failed to register user")
    }

	return http.StatusOK, nil
}
