package helpers

import (
	"os"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte
var once sync.Once

// GetJWTSecret ensures the secret is loaded once after .env is loaded
func GetJWTSecret() []byte {
	once.Do(func() {
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			panic("JWT_SECRET environment variable not set")
		}
		jwtSecret = []byte(secret)
	})
	return jwtSecret
}

// GenerateToken creates a JWT token with the email claim and 72h expiry
func GenerateToken(email string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(72 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(GetJWTSecret()) // <-- use exported function
}
