package client_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var token string

func TestMain(m *testing.M) {
	secret := []byte("secret")
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Minute * 2).Unix(),
	}
	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	var err error
	token, err = tokenObj.SignedString(secret)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(m.Run())
}
