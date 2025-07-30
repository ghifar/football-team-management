package usecases

import (
	"errors"
	"football-team-management/internal/domain/user"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	GenerateToken(username, password string) (string, error)
	ValidateToken(tokenString string) (*user.Claims, error)
}

type authService struct {
	jwtSecret []byte
	getUser   GetUser
}

func NewAuthService(jwtSecret string, getUser GetUser) AuthService {
	return &authService{
		jwtSecret: []byte(jwtSecret),
		getUser:   getUser,
	}
}

func (a *authService) GenerateToken(username, password string) (string, error) {
	// Get user from repository
	userData, err := a.getUser.Execute(nil, username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	// Create claims
	claims := &user.Claims{
		Username: userData.Username,
		Role:     userData.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 24 hours
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(a.jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (a *authService) ValidateToken(tokenString string) (*user.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &user.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return a.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*user.Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
