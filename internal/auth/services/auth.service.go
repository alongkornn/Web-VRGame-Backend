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
	"github.com/alongkornn/Web-VRGame-Backend/pkg/utils"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
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
		"user_id":  user.ID,
		"email":    user.Email,
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

func SendVerificationEmail(ctx context.Context, e string) (int, error) {
	email := e
	if email == "" {
		return http.StatusBadRequest, errors.New("email is required")
	}

	// สร้างลิงก์ยืนยันอีเมล
	client, err := utils.FirebaseApp.Auth(context.Background())
	if err != nil {
		return http.StatusInternalServerError, errors.New("failed to get Firebase Auth client")
	}

	link, err := client.EmailVerificationLink(context.Background(), email)
	if err != nil {
		return http.StatusInternalServerError, errors.New("failed to generate email verification link")
	}

	// ส่งอีเมลยืนยัน
	if err := sendEmail(email, link); err != nil {
		return http.StatusInternalServerError, errors.New("failed to send email")
	}

	return http.StatusOK, nil
}

func sendEmail(to, link string) error {
	from := "alongkornp5363@gmail.com" // อีเมลผู้ส่ง
	password := "djtwoggiuvoiswot"     // รหัสผ่านสำหรับ SMTP

	// สร้างอีเมล
	msg := gomail.NewMessage()
	msg.SetHeader("From", from)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", "Email Verification")
	msg.SetBody("text/html", fmt.Sprintf(`
        <p>Thank you for registering! Please verify your email by clicking the link below:</p>
        <a href="%s">Verify Email</a>
    `, link))

	// ส่งอีเมลผ่าน SMTP
	dialer := gomail.NewDialer("smtp.gmail.com", 587, from, password)

	return dialer.DialAndSend(msg)
}
