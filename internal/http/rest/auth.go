package rest

import (
	"errors"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

// authenticationRequest specifies the structure of json when authenticating
type authenticationRequest struct {
	Password string `json:"password"`
}

// errSigningJwt represents the error returned when Token can not be signed
var errSigningJwt error = errors.New("Could not sign JWT Token")

// generateTokenPair returns signed access & refresh tokens
func generateTokenPair() (string, string, error) {
	secret := JwtSecret()
	// Access
	access := jwt.New(jwt.SigningMethodHS256)
	accessClaims := access.Claims.(jwt.MapClaims)
	accessClaims["name"] = "El Hamza"
	accessClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	signedAccess, err := access.SignedString(secret)
	if err != nil {
		return "", "", err
	}
	// Refresh
	refresh := jwt.New(jwt.SigningMethodHS256)
	refreshClaims := refresh.Claims.(jwt.MapClaims)
	refreshClaims["exp"] = time.Now().Add(time.Hour * 8).Unix()
	signedRefresh, err := refresh.SignedString(secret)
	if err != nil {
		return "", "", err
	}
	return signedAccess, signedRefresh, nil
}

// Login handler authenticates user and returns a JWT Token
func (h *Handler) Login(c echo.Context) error {
	logger := log.WithFields(log.Fields{
		"handler":   "Login",
		"remote_ip": c.RealIP(),
	})
	// Unmarshal JSON
	var req authenticationRequest
	if err := c.Bind(&req); err != nil {
		var (
			msg     string = errInvalidJSON.Error()
			details string = httpErrorMsg(err)
			code    int    = errToHTTPCode(errInvalidJSON, "auth")
		)
		logger.Error(msg + " | " + details)
		return c.String(code, msg)
	}
	// Authenticate
	if err := h.authenticator.Authenticate(req.Password); err != nil {
		msg := err.Error()
		logger.Info(msg)
		code := errToHTTPCode(err, "auth")
		return c.String(code, msg)
	}
	// Sign and return Access and Refresh Tokens
	access, refresh, errSign := generateTokenPair()
	if errSign != nil {
		var (
			msg  string = errSigningJwt.Error()
			code int    = errToHTTPCode(errSigningJwt, "auth")
		)
		logger.Error(msg + " : " + errSign.Error())
		return c.String(code, msg)
	}
	body := map[string]string{
		"at": access,
		"rt": refresh,
	}
	return c.JSON(http.StatusOK, body)
}
