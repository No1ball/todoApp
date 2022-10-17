package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/No1ball/todo-app/internal/models"
	"github.com/No1ball/todo-app/internal/repository"
	"github.com/dgrijalva/jwt-go"
	"time"
)

const (
	salt      = "hjklfkjsadnflaskdf"
	signInKey = "ghjklfsdfffdsfsfsw11f"
	tokenTTL  = 12 * time.Hour
)

type AuthService struct {
	repo repository.Authorization
}
type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user models.User) (int, error) {
	user.Password = generatePassword(user.Password)
	return s.repo.CreateUser(user)
}

func (s *AuthService) GenerateToken(user models.SignInInput) (string, error) {
	var userWithPassword models.SignInInput
	userWithPassword.Username = user.Username
	userWithPassword.Password = generatePassword(user.Password)
	getUser, err := s.repo.GetUser(userWithPassword)

	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		getUser.Id,
	})
	return token.SignedString([]byte(signInKey))
}

func (s *AuthService) ParseToken(token string) (int, error) {
	newToken, err := jwt.ParseWithClaims(token, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(signInKey), nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := newToken.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}
	return claims.UserId, nil
}
func generatePassword(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
