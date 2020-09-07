package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Login handler authenticates user and returns a JWT Token
func (h *Handler) Login(c echo.Context) error {
	var authReq struct {
		Password string `json:"password"`
	}
	if err := c.Bind(&authReq); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	if authReq.Password == "mytestpass" {
		resp := map[string]string{
			"at": "some access token not yet implemented",
		}
		return c.JSON(http.StatusOK, resp)
	}
	return c.String(http.StatusUnauthorized, "")
}
