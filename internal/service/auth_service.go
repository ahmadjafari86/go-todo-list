package service

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/ahmadjafari86/go-todo-list/internal/models"
	"github.com/ahmadjafari86/go-todo-list/internal/repository"
)

type AuthService interface {
	Register(email, password string) (*models.User, error)
	Login(email, password string) (string, error)
	ParseToken(tokenStr string) (*jwt.RegisteredClaims, error)
}

type authService struct {
	users         repository.UserRepository
	jwtSecret     string
	jwtExpMinutes int
}

func NewAuthService(users repository.UserRepository) AuthService {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "change_this_in_prod"
	}
	exp := 60
	if v := os.Getenv("JWT_EXP_MINUTES"); v != "" {
		fmt.Sscanf(v, "%d", &exp)
	}
	return &authService{users: users, jwtSecret: secret, jwtExpMinutes: exp}
}

func (s *authService) Register(email, password string) (*models.User, error) {
	ex, err := s.users.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	if ex != nil {
		return nil, errors.New("email already registered")
	}
	hpw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	u := &models.User{Email: email, PasswordHash: string(hpw)}
	if err := s.users.Create(u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *authService) Login(email, password string) (string, error) {
	u, err := s.users.GetByEmail(email)
	if err != nil {
		return "", err
	}
	if u == nil {
		return "", errors.New("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	exp := time.Now().Add(time.Duration(s.jwtExpMinutes) * time.Minute)
	claims := jwt.RegisteredClaims{
		Subject:   fmt.Sprint(u.ID),
		ExpiresAt: jwt.NewNumericDate(exp),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokStr, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}
	return tokStr, nil
}

func (s *authService) ParseToken(tokenStr string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
