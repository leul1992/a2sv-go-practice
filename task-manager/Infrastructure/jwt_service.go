package Infrastructure

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JWTService interface {
	GenerateToken(userID primitive.ObjectID, username, role string) (string, error)
}

type jwtService struct{}

func NewJWTService() JWTService {
	return &jwtService{}
}

func (j *jwtService) GenerateToken(userID primitive.ObjectID, username, role string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", ErrJWTConfig
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":      userID.Hex(),
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(72 * time.Hour).Unix(),
	})

	return jwtToken.SignedString([]byte(secret))
}

var ErrJWTConfig = fmt.Errorf("JWT secret not set in environment")
