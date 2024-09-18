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
	"google.golang.org/api/iterator"
)

// user

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
		fmt.Printf("Error adding document: %v\n", err)
		return http.StatusInternalServerError, errors.New("failed to register user")
	}

	return http.StatusOK, nil
}

func Login(email, password string, ctx context.Context) (*dto.ResponseLogin, int, error) {
	hasUser := config.DB.Collection("User").Where("email", "==", email).Limit(1)
	doc, err := hasUser.Documents(ctx).Next()
	if err != nil {
		return nil, http.StatusBadRequest, errors.New("user not found")
	}

	var user models.User
	doc.DataTo(&user)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, http.StatusUnauthorized, errors.New("invalid password")
	}

	token, err := generateToken(&user)
	if err != nil {
		return nil, http.StatusUnauthorized, errors.New("failed to create token")
	}

	data := dto.ResponseLogin{
		Token: token,
	}

	return &data, http.StatusOK, nil
}

func GetUser(ctx context.Context) ([]*models.User, int, error) {
	iter := config.DB.Collection("User").Where("is_deleted", "==", false).Documents(ctx)
	defer iter.Stop()

	var users []*models.User
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, http.StatusInternalServerError, errors.New("somethig went wrong")
		}

		var user models.User
		err = doc.DataTo(&user)
		if err != nil {
			return nil, http.StatusInternalServerError, errors.New("somethig went wrong")
		}

		users = append(users, &user)
	}

	if len(users) <= 0 {
		return nil, http.StatusOK, errors.New("users is empty")
	}

	return users, http.StatusOK, nil
}

func generateToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.GetEnv("jwt.secret_key")))
}