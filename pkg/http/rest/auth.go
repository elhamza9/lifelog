package rest

import (
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// Login handler authenticates user and returns a JWT Token
func (h *Handler) Login(c echo.Context) error {
	// Unmarshal JSON
	var authReq struct {
		Password string `json:"password"`
	}
	if err := c.Bind(&authReq); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	const passwordMinLength int = 8
	if len(authReq.Password) < passwordMinLength {
		return c.String(http.StatusBadRequest, "Password too short.")
	}
	// Get Correct Password
	hash := os.Getenv("LFLG_PASS_HASH")
	if hash == "" {
		return c.String(http.StatusInternalServerError, "Can not perform authentication because no Password was found in system. ")
	}
	// Authenticate
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(authReq.Password)); err != nil {
		return c.String(http.StatusUnauthorized, "Wrong Password")
	}
	// Return JWT Token
	secret := JwtSecret()
	token := jwt.New(jwt.SigningMethodHS256)
	signed, err := token.SignedString(secret)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	body := map[string]string{
		"at": signed,
	}
	return c.JSON(http.StatusOK, body)

}
