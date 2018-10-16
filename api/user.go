package api

import (
	"github.com/labstack/echo"
	er "github.com/pintobikez/popmeet/errors"
	repo "github.com/pintobikez/popmeet/repository"
	"net/http"
	"strconv"
)

type UserApi struct {
	rp repo.Repository
}

func (a *UserApi) SetRepository(rpo repo.Repository) {
	a.rp = rpo
}

// Handler to GET Interest
func (a *UserApi) GetUser() echo.HandlerFunc {
	return func(c echo.Context) error {

		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			return c.JSON(http.StatusBadRequest, &er.ErrResponse{er.ErrContent{http.StatusBadRequest, err.Error()}})
		}

		// Find the user
		resp, err := a.rp.FindUserById(id)
		if err != nil {
			return c.JSON(http.StatusNotFound, &er.ErrResponse{er.ErrContent{er.ErrorUserNotFound, err.Error()}})
		}
		// Find the user profile
		resp.Profile, err = a.rp.FindProfileByUserId(resp.id)
		if err != nil {
			return c.JSON(http.StatusNotFound, &er.ErrResponse{er.ErrContent{er.ErrorUserProfileNotFound, err.Error()}})
		}
		// Find the user security
		resp.Profile, err = a.rp.FindSecurityInfoByUserId(resp.id)
		if err != nil {
			return c.JSON(http.StatusNotFound, &er.ErrResponse{er.ErrContent{er.ErrorUserProfileNotFound, err.Error()}})
		}

		return c.JSON(http.StatusOK, resp)
	}
}
