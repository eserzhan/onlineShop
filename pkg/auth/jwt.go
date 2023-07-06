package auth

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// TokenManager provides logic for JWT & Refresh tokens generation and parsing.
type TokenManager interface {
	NewJWT(userId, role string, ttl time.Duration) (string, error)
	Parse(accessToken string) (map[string]interface{}, error)
	NewRefreshToken() (string, error)
}

type Manager struct {
	signingKey string
}

func NewManager(signingKey string) (*Manager, error) {
	if signingKey == "" {
		return nil, errors.New("empty signing key")
	}

	return &Manager{signingKey: signingKey}, nil
}

type CustomClaims struct {
	jwt.StandardClaims
	Role string `json:"role"`
}

func (m *Manager) NewJWT(userId, role string, ttl time.Duration) (string, error) {
	claims := CustomClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(ttl).Unix(),
			Subject:   userId,
		},
		Role: role, 
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(m.signingKey))
}

func (m *Manager) Parse(accessToken string) (map[string]interface{}, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(m.signingKey), nil
	})
	if err != nil {
		return map[string]interface{}{}, err
	}
	
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return map[string]interface{}{}, fmt.Errorf("error get user claims from token")
	}

	userID := claims["sub"].(string)
	role := claims["role"].(string)

	userInfo := map[string]interface{}{
		"userID": userID,
		"role":   role,
	}

	return userInfo, nil

}

func (m *Manager) NewRefreshToken() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}
