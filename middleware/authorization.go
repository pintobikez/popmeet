package middleware

import (
	"github.com/labstack/echo"
	er "github.com/pintobikez/popmeet/errors"
	"github.com/pintobikez/popmeet/secure"
	"net/http"
)

// Authorization Middleware
func Authorization(sec *secure.TokenManager) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			claims, err := sec.ValidateToken(c.Request().Header.Get(echo.HeaderAuthorization), "")
			if err != nil {
				return c.JSON(http.StatusUnauthorized, er.GeneralErrorJson(http.StatusUnauthorized, "Invalid token"))
			}

			c.Set("claims", claims)

			return next(c)
		}
	}
}
