package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// loginRequest specifies the structure of json when authenticating
type loginRequest struct {
	Password string `json:"password"`
}

// refreshRequest specifies the structure of json when refreshing Tokens
type refreshRequest struct {
	RefreshToken string `json:"refresh"`
}

// errSigningJwt represents the error returned when Token can not be signed
var errSigningJwt error = errors.New("Could not sign JWT Token")

// jwtSignErrHandler accepts the error returned when signing a JWT Token
// and returns:
//	- Http Code to be returned to the client
//	- Msg body to be returned to the client
//	- Log message to be logged
func jwtSignErrHandler(err error) (code int, respMsg string, logMsg string) {
	code = errToHTTPCode(errSigningJwt, "auth")
	respMsg = errSigningJwt.Error()
	logMsg = respMsg + " : " + err.Error()
	return code, respMsg, logMsg
}

// Login handler authenticates user and returns a JWT Token
func (h *Handler) Login(c echo.Context) error {
	// Unmarshal JSON
	var req loginRequest
	if err := c.Bind(&req); err != nil {
		var (
			msg     string = errInvalidJSON.Error()
			details string = httpErrorMsg(err)
			code    int    = errToHTTPCode(errInvalidJSON, "auth")
		)
		logrus.Error(msg + " | " + details)
		return c.String(code, msg)
	}
	// Authenticate
	if err := h.authenticator.Authenticate(req.Password); err != nil {
		msg := err.Error()
		logrus.Error(msg)
		code := errToHTTPCode(err, "auth")
		return c.String(code, msg)
	}
	logrus.Info("Authentication successful")
	// Generate and return Access/Refresh Tokens
	access, err := generateAccessToken()
	if err != nil {
		code, msg, logMsg := jwtSignErrHandler(err)
		logrus.Error(logMsg)
		return c.String(code, msg)
	}
	logrus.Info("Generated Access Token")
	refresh, err := generateRefreshToken()
	if err != nil {
		code, msg, logMsg := jwtSignErrHandler(err)
		logrus.Error(logMsg)
		return c.String(code, msg)
	}
	logrus.Info("Generated Refresh Token")
	body := map[string]string{
		"at": access,
		"rt": refresh,
	}
	return c.JSON(http.StatusOK, body)
}

// RefreshToken handler accepts a refresh token
// and returns a new access/refresh token pair
func (h *Handler) RefreshToken(c echo.Context) error {
	// Unmarshal JSON
	var req refreshRequest
	if err := c.Bind(&req); err != nil {
		var (
			msg     string = errInvalidJSON.Error()
			details string = httpErrorMsg(err)
			code    int    = errToHTTPCode(errInvalidJSON, "auth")
		)
		logrus.Error(msg + " | " + details)
		return c.String(code, msg)
	}
	if req.RefreshToken == "" {
		code := errToHTTPCode(errInvalidJSON, "auth")
		msg := "No Refresh Token Provided"
		logrus.Error(msg)
		return c.String(code, msg)
	}
	logrus.Info("Extracted refresh token successfully")
	// Parse Token
	_, err := jwt.Parse(req.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			msg := fmt.Sprintf("Unexpected signing method: %v", token.Header["alg"])
			logrus.Error(msg)
			return nil, errors.New(msg)
		}
		return jwtRefreshSecret(), nil
	})
	if err != nil {
		msg := "Refresh Token is Invalid"
		logrus.Error(msg + " = " + err.Error())
		return c.String(http.StatusUnprocessableEntity, msg)
	}
	logrus.Info("Refresh Token validation successful")
	// Generate and return new Access Token
	access, err := generateAccessToken()
	if err != nil {
		code, msg, logMsg := jwtSignErrHandler(err)
		logrus.Error(logMsg)
		return c.String(code, msg)
	}
	logrus.Info("Generated Access Token")
	body := map[string]string{
		"at": access,
	}
	return c.JSON(http.StatusOK, body)
}
