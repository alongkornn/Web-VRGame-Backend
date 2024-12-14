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
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// ลงทะเบียน
func Register(ctx context.Context, registerDTO *dto.RegisterDTO) (int, error) {
	hasUser := config.DB.Collection("User").
		Where("email", "==", registerDTO.Email).
		Limit(1)

	userDoc, err := hasUser.Documents(ctx).GetAll()
	if err != nil {
		return http.StatusInternalServerError, errors.New("error checking user existence")
	}

	if len(userDoc) > 0 {
		return http.StatusConflict, errors.New("email already registered")
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(registerDTO.Password), bcrypt.DefaultCost)
	if err != nil {
		return http.StatusBadRequest, errors.New("hash password is error")
	}
	user := models.User{
		ID:                   uuid.New().String(),
		FirstName:            registerDTO.FirstName,
		LastName:             registerDTO.LastName,
		Email:                registerDTO.Email,
		Password:             string(hashPassword),
		Class:                registerDTO.Class,
		Number:               registerDTO.Number,
		Level:                1,
		Score:                0,
		CurrentCheckpoint:    nil,
		CompletedCheckpoints: nil,
		Role:                 models.Player,
		Status:               models.Pending,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
		Is_Deleted:           false,
	}

	_, _, err = config.DB.Collection("User").Add(ctx, user)
	if err != nil {
		fmt.Printf("Error adding document: %v\n", err)
		return http.StatusInternalServerError, errors.New("failed to register user")
	}

	return http.StatusOK, nil
}

// เข้าสู่ระบบ
func Login(email, password string, ctx context.Context) (string, int, error) {
	hasUser := config.DB.Collection("User").Where("email", "==", email).Limit(1)

	userDoc, err := hasUser.Documents(ctx).GetAll()
	if err != nil || len(userDoc) == 0 {
		return "", http.StatusBadRequest, errors.New("user not found")
	}

	var user models.User
	if err := userDoc[0].DataTo(&user); err != nil {
		return "nil", http.StatusInternalServerError, errors.New("error retrieving user data")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "nil", http.StatusUnauthorized, errors.New("invalid password")
	}

	token, err := generateToken(&user)
	if err != nil {
		return "nil", http.StatusUnauthorized, errors.New("failed to create token")
	}

	return token, http.StatusOK, nil
}

// สร้าง token ขึ้นมา
func generateToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"username": user.FirstName,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.GetEnv("jwt.secret_key")))

	if err != nil {
		return "", errors.New("invalid create token")
	}
	return tokenString, nil
}
