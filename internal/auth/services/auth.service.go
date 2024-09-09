package services

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/alongkornn/Web-VRGame-Backend/internal/auth/dto"
	"github.com/alongkornn/Web-VRGame-Backend/internal/auth/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)


func Register(userRegister dto.RegisterDTO) (int, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(userRegister.Password), 10)
	if err != nil {
		return http.StatusBadRequest, errors.New("hash password is error")
	}
	 user := models.User{
		ID: uuid.New().String(),
		FirstName: userRegister.FirstName,
		LastName: userRegister.LastName,
		Email: userRegister.Email,
		Password: string(hashPassword),
		Class: userRegister.Class,
		Number: userRegister.Number,
		Role: models.Player,
		Score: 0,
		Status: models.Pending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Is_Deleted: false,
	}

	fmt.Println(user)


	return http.StatusOK, nil
}