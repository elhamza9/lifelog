package server

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// jwtAccessSecret returns Jwt Signing Key for Access Tokens as a bytes slice.
func jwtAccessSecret() []byte {
	secret := os.Getenv("LFLG_JWT_ACCESS_SECRET")
	return []byte(secret)
}

// jwtRefreshSecret returns Jwt Signing Key for Refresh Tokens as a bytes slice.
func jwtRefreshSecret() []byte {
	secret := os.Getenv("LFLG_JWT_REFRESH_SECRET")
	return []byte(secret)
}

const (
	// accessTokenExpDuration represents how much time access token lives before it must be changed.
	accessTokenExpDuration time.Duration = time.Duration(time.Minute * 15)
	// refreshTokenExpDuration represents how much time refresh token lives before it must be changed.
	refreshTokenExpDuration time.Duration = time.Duration(time.Hour * 6)
)

// accessTokenClaims represents claims used in Access Token.
type accessTokenClaims struct {
	Name string `json:"name"`
	jwt.StandardClaims
}

// generateAccessToken generates & returns a signed access token.
func generateAccessToken() (string, error) {
	secret := jwtAccessSecret()
	now := time.Now()
	claims := &accessTokenClaims{
		"El Hamza",
		jwt.StandardClaims{
			ExpiresAt: now.Add(accessTokenExpDuration).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return signed, nil
}

// generateTokenPair generates & returns a signed refresh token.
func generateRefreshToken() (string, error) {
	secret := jwtRefreshSecret()
	now := time.Now()
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = now.Add(refreshTokenExpDuration).Unix()
	signed, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return signed, nil
}
