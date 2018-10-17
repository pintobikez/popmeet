package api

import (
	"github.com/labstack/echo"
	"github.com/pintobikez/popmeet/api/models"
	er "github.com/pintobikez/popmeet/errors"
	repo "github.com/pintobikez/popmeet/repository"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"strconv"
)

type UserApi struct {
	rp       repo.Repository
	validate *validator.Validate
}

func (a *UserApi) New(rpo repo.Repository) {
	a.rp = rpo
	a.validate = validator.New()
}

func (a *UserApi) SetRepository(rpo repo.Repository) {
	a.rp = rpo
}

// Handler to GET Interest
func (a *UserApi) GetUser() echo.HandlerFunc {
	return func(c echo.Context) error {

		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, &er.ErrResponse{er.ErrContent{http.StatusBadRequest, err.Error()}})
		}

		// Find the user
		resp, err := a.rp.FindUserById(id)
		if err != nil {
			return c.JSON(http.StatusNotFound, &er.ErrResponse{er.ErrContent{er.ErrorUserNotFound, err.Error()}})
		}
		// Find the user profile
		resp.Profile, err = a.rp.FindUserProfileByUserId(resp.ID)
		if err != nil {
			return c.JSON(http.StatusNotFound, &er.ErrResponse{er.ErrContent{er.ErrorUserProfileNotFound, err.Error()}})
		}
		// Find the user security
		resp.Security, err = a.rp.FindSecurityInfoByUserId(resp.ID)
		if err != nil {
			return c.JSON(http.StatusNotFound, &er.ErrResponse{er.ErrContent{er.ErrorUserProfileNotFound, err.Error()}})
		}

		return c.JSON(http.StatusOK, resp)
	}
}

// Handler to PUT User
func (a *UserApi) PutUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		var err error

		u := new(models.NewUser)
		if err := c.Bind(u); err != nil {
			return c.JSON(http.StatusBadRequest, &er.ErrResponse{er.ErrContent{http.StatusBadRequest, err.Error()}})
		}

		if err := a.validate.Struct(u); err != nil {
			return c.JSON(http.StatusBadRequest, &er.ErrResponse{er.ErrContent{http.StatusBadRequest, err.Error()}})
		}

		if u.Password == "" && u.Provider == 0 {
			return c.JSON(http.StatusBadRequest, &er.ErrResponse{er.ErrContent{http.StatusBadRequest, "A password or a Login Provider must be provided"}})
		}

		ur := &models.User{Name: u.Name, Email: u.Email, Active: true, Security: &models.UserSecurity{LastMachine: c.RealIP()}}

		// Hash the password
		if u.Password != "" {
			ur.Security.Hash, err = a.hashPassword(u.Password)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, &er.ErrResponse{er.ErrContent{http.StatusInternalServerError, err.Error()}})
			}
		} else { //find the login provider
			ur.Security.Provider, err = a.rp.FindLoginProviderById(u.Provider)
			if err != nil {
				if ur.Security.Provider.ID == -1 {
					return c.JSON(http.StatusBadRequest, &er.ErrResponse{er.ErrContent{http.StatusBadRequest, err.Error()}})
				} else {
					return c.JSON(http.StatusInternalServerError, &er.ErrResponse{er.ErrContent{http.StatusInternalServerError, err.Error()}})
				}
			}
		}

		err = a.rp.InsertUser(ur)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, &er.ErrResponse{er.ErrContent{http.StatusInternalServerError, err.Error()}})
		}

		return c.JSON(http.StatusOK, ur)
	}
}

func (a *UserApi) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (a *UserApi) checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Handler to POST User
/*func (a *UserApi) PostUser() echo.HandlerFunc {
	return func(c echo.Context) error {

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, &er.ErrResponse{er.ErrContent{http.StatusBadRequest, err.Error()}})
		}

		u := new(models.User)
		if err = c.Bind(u); err != nil {
			return c.JSON(http.StatusBadRequest, &er.ErrResponse{er.ErrContent{http.StatusBadRequest, err.Error()}})
		}
		if err = a.validate.Struct(u); err != nil {
			return c.JSON(http.StatusBadRequest, &er.ErrResponse{er.ErrContent{http.StatusBadRequest, err.Error()}})
		}

		return c.JSON(http.StatusOK, u)
	}
}
*/
