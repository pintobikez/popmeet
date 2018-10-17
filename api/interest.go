package api

import (
	"github.com/labstack/echo"
	er "github.com/pintobikez/popmeet/errors"
	repo "github.com/pintobikez/popmeet/repository"
	"net/http"
	"strconv"
)

type InterestApi struct {
	rp repo.Repository
}

func (a *InterestApi) New(rpo repo.Repository) {
	a.rp = rpo
}

func (a *InterestApi) SetRepository(rpo repo.Repository) {
	a.rp = rpo
}

// Handler to GET Interest
func (a *InterestApi) GetInterest() echo.HandlerFunc {
	return func(c echo.Context) error {

		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, &er.ErrResponse{er.ErrContent{http.StatusBadRequest, err.Error()}})
		}

		resp, err := a.rp.FindInterestById(id)

		if err != nil {
			return c.JSON(http.StatusNotFound, &er.ErrResponse{er.ErrContent{er.ErrorInterestNotFound, err.Error()}})
		}

		return c.JSON(http.StatusOK, resp)
	}
}

// Handler to GET All Interests
func (a *InterestApi) GetAllInterest() echo.HandlerFunc {
	return func(c echo.Context) error {

		resp, err := a.rp.GetAllInterests()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, &er.ErrResponse{er.ErrContent{er.ErrorInterestsNotFound, err.Error()}})
		}

		return c.JSON(http.StatusOK, resp)
	}
}
