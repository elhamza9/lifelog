package rest

import (
	"errors"
	"net/http"
	"os"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// authenticationRequest specifies the structure of json when authenticating
type authenticationRequest struct {
	Password string `json:"password"`
}

// passwordMinLength specifies minimum length given password should be.
const passwordMinLength int = 8
const passwordMaxLength int = 256

// getPasswordHash returns bcrypt hash of correct password
func getPasswordHash() string {
	return os.Getenv("LFLG_PASS_HASH")
}

// Errors
var (
	errPasswordLength       error = errors.New("Password must be " + strconv.Itoa(passwordMinLength) + " ~ " + strconv.Itoa(passwordMaxLength) + " characters")
	errIncorrectCredentials error = errors.New("Incorrect credentials")
	errHashNotFound         error = errors.New("Could not perform authentication")
	errSigningJwt           error = errors.New("Could not sign JWT Token")
)

// validatePassword validates password format
func validatePassword(pass string) error {
	if len(pass) < passwordMinLength || len(pass) > passwordMaxLength {
		return errPasswordLength
	}
	return nil
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
	// Check Password Length
	if err := validatePassword(req.Password); err != nil {
		var (
			msg     string = err.Error()
			details string = strconv.Itoa(len(req.Password)) + " characters"
			code    int    = errToHTTPCode(err, "auth")
		)
		logger.Info(msg + " : " + details)
		return c.String(code, msg)
	}
	// Get Correct Password
	hash := getPasswordHash()
	if hash == "" {
		msg := errHashNotFound.Error()
		logger.Error(msg + " : " + "No Password Hash was found in system.")
		return c.String(http.StatusInternalServerError, msg)
	}
	// Authenticate
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)); err != nil {
		msg := errIncorrectCredentials.Error()
		logger.Info(msg)
		code := errToHTTPCode(errIncorrectCredentials, "auth")
		return c.String(code, msg)
	}
	// Return JWT Token
	secret := JwtSecret()
	token := jwt.New(jwt.SigningMethodHS256)
	signed, err := token.SignedString(secret)
	if err != nil {
		msg := errSigningJwt.Error()
		logger.Error(msg + " : " + err.Error())
		code := errToHTTPCode(errSigningJwt, "auth")
		return c.String(code, msg)
	}
	body := map[string]string{
		"at": signed,
	}
	return c.JSON(http.StatusOK, body)
}
