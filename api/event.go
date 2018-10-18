package api

import (
	"github.com/labstack/echo"
	"github.com/pintobikez/popmeet/api/models"
	er "github.com/pintobikez/popmeet/errors"
	repo "github.com/pintobikez/popmeet/repository"
	stru "github.com/pintobikez/popmeet/secure/structures"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"strconv"
)

type EventApi struct {
	rp       repo.Repository
	validate *validator.Validate
}

func (a *EventApi) New(rpo repo.Repository) {
	a.rp = rpo
	a.validate = validator.New()
}

func (a *EventApi) SetRepository(rpo repo.Repository) {
	a.rp = rpo
}

// GetEvent Handler to GET Event
func (a *EventApi) GetEvent() echo.HandlerFunc {
	return func(c echo.Context) error {

		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, er.GeneralErrorJson(http.StatusBadRequest, err.Error()))
		}

		resp, err := a.rp.GetEventById(id)
		if err != nil {
			return c.JSON(http.StatusNotFound, er.GeneralErrorJson(er.ErrorEventNotFound, err.Error()))
		}

		return c.JSON(http.StatusOK, resp)
	}
}

// PutEvent Handler to PUT Event
func (a *EventApi) PutEvent() echo.HandlerFunc {
	return func(c echo.Context) error {

		u := new(models.NewEvent)
		if err := c.Bind(u); err != nil {
			return c.JSON(http.StatusBadRequest, er.GeneralErrorJson(http.StatusBadRequest, err.Error()))
		}

		if err := a.validate.Struct(u); err != nil {
			return c.JSON(http.StatusUnprocessableEntity, er.ValidationErrorJson(http.StatusUnprocessableEntity, err))
		}

		cl := c.Get("claims").(*stru.TokenClaims)
		// Get the user by the claim ID
		ur, err := a.rp.GetUserById(cl.ID)
		if err != nil {
			return c.JSON(http.StatusNotFound, er.GeneralErrorJson(er.ErrorUserNotFound, err.Error()))
		}

		ev := &models.Event{StartDate: u.StartDate, EndDate: u.EndDate, Location: u.Location, Longitude: u.Longitude, Latitude: u.Latitude, Active: u.Active, CreatedBy: ur}

		//Save the event
		err = a.rp.InsertEvent(ev)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, er.GeneralErrorJson(http.StatusInternalServerError, err.Error()))
		}

		//Get the complete info from the event to return it
		ev, err = a.rp.GetEventById(ev.ID)
		if err != nil {
			return c.JSON(http.StatusNotFound, er.GeneralErrorJson(er.ErrorEventNotFound, err.Error()))
		}

		return c.JSON(http.StatusOK, ev)
	}
}

// AddUserToEvent Handler to PUT a User in an Event
func (a *EventApi) AddUserToEvent() echo.HandlerFunc {
	return func(c echo.Context) error {

		// gets the event id
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, er.GeneralErrorJson(http.StatusBadRequest, err.Error()))
		}

		// Get the user by the claim ID
		cl := c.Get("claims").(*stru.TokenClaims)

		//Check if the event exists and its active
		ex, err := a.rp.FindEventById(id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, er.GeneralErrorJson(http.StatusInternalServerError, err.Error()))
		}
		if !ex {
			return c.JSON(http.StatusNotFound, er.GeneralErrorJson(er.ErrorEventNotFound, "Event doesn't exist"))
		}

		//Check if the user exists and its active
		ex, err = a.rp.FindUserById(cl.ID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, er.GeneralErrorJson(http.StatusInternalServerError, err.Error()))
		}
		if !ex {
			return c.JSON(http.StatusNotFound, er.GeneralErrorJson(er.ErrorUserNotFound, "User doesn't exist"))
		}

		//Add the user to the event
		err = a.rp.AddUserToEvent(id, cl.ID)
		if err != nil {
			if err.Error() == strconv.Itoa(er.ErrorCantAddUSerToEvent) {
				return c.JSON(http.StatusBadRequest, er.GeneralErrorJson(er.ErrorCantAddUSerToEvent, "Can't add creator as user"))
			}
			return c.JSON(http.StatusInternalServerError, er.GeneralErrorJson(http.StatusInternalServerError, err.Error()))
		}

		return c.NoContent(http.StatusOK)
	}
}

// RemoveUserFromEvent Handler to DELETE a User from an Event
func (a *EventApi) RemoveUserFromEvent() echo.HandlerFunc {
	return func(c echo.Context) error {

		// gets the event id
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, er.GeneralErrorJson(http.StatusBadRequest, err.Error()))
		}

		// Get the user by the claim ID
		cl := c.Get("claims").(*stru.TokenClaims)

		//Check if the event exists and its active
		ex, err := a.rp.FindEventById(id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, er.GeneralErrorJson(http.StatusInternalServerError, err.Error()))
		}
		if !ex {
			return c.JSON(http.StatusNotFound, er.GeneralErrorJson(er.ErrorEventNotFound, "Event doesn't exist"))
		}

		//Check if the user exists and its active
		ex, err = a.rp.FindUserById(cl.ID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, er.GeneralErrorJson(http.StatusInternalServerError, err.Error()))
		}
		if !ex {
			return c.JSON(http.StatusNotFound, er.GeneralErrorJson(er.ErrorUserNotFound, "User doesn't exist"))
		}

		//Remove the user from the event
		err = a.rp.RemoveUserFromEvent(id, cl.ID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, er.GeneralErrorJson(http.StatusInternalServerError, err.Error()))
		}

		return c.NoContent(http.StatusOK)
	}
}

// FindEvents Handler to find Events by given search parameters
func (a *EventApi) FindEvents() echo.HandlerFunc {
	return func(c echo.Context) error {

		return c.NoContent(http.StatusOK)
	}
}
