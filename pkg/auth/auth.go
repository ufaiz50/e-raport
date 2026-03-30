package auth

import (
	"crypto/rand"
	"encoding/base64"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// Claims struct to be encoded to JWT
type Claims struct {
	Username string  `json:"username"`
	Role     string  `json:"role"`
	SchoolID *string `json:"school_id,omitempty"`
	jwt.StandardClaims
}

var JwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

func refreshJwtKey() []byte {
	if key := os.Getenv("REFRESH_JWT_SECRET_KEY"); key != "" {
		return []byte(key)
	}
	return JwtKey
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func GenerateAccessToken(username string, role string, schoolID *string) (string, error) {
	expirationTime := time.Now().Add(15 * time.Minute).Unix()

	claims := &Claims{
		Username: username,
		Role:     role,
		SchoolID: schoolID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime,
			Issuer:    username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GenerateRefreshToken(username string, role string, schoolID *string) (string, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour).Unix()

	claims := &Claims{
		Username: username,
		Role:     role,
		SchoolID: schoolID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime,
			Issuer:    username,
			Subject:   "refresh",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(refreshJwtKey())
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseRefreshToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return refreshJwtKey(), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}
	if claims.Subject != "refresh" {
		return nil, jwt.ErrSignatureInvalid
	}
	return claims, nil
}

// Backward compatibility
func GenerateToken(username string, role string, schoolID *string) (string, error) {
	return GenerateAccessToken(username, role, schoolID)
}

func GenerateRandomKey() string {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		panic("Failed to generate random key: " + err.Error())
	}

	return base64.StdEncoding.EncodeToString(key)
}
