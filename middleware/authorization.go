package middleware

import (
	"github.com/labstack/echo"
	"github.com/pintobikez/popmeet/api/models"
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
				return c.JSON(http.StatusUnauthorized, &er.ErrResponse{er.ErrContent{http.StatusUnauthorized, "Invalid token"}})
			}

			//lets authorize user actions
			if c.Request().Method == "POST" && c.Request().RequestURI == "/user" {
				u := new(models.User)
				if err := c.Bind(u); err != nil {
					return c.JSON(http.StatusBadRequest, &er.ErrResponse{er.ErrContent{http.StatusBadRequest, err.Error()}})
				}

				if claims.ID != u.ID {
					return c.JSON(http.StatusUnauthorized, &er.ErrResponse{er.ErrContent{http.StatusUnauthorized, "Not authorized to update this user"}})
				}
			}

			c.Set("claims", claims)

			return next(c)
		}
	}
}
