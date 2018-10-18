package api

import (
	"github.com/labstack/echo"
	"github.com/pintobikez/popmeet/api/models"
	er "github.com/pintobikez/popmeet/errors"
	repo "github.com/pintobikez/popmeet/repository"
	"github.com/pintobikez/popmeet/secure"
	tok "github.com/pintobikez/popmeet/secure/structures"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"strconv"
)

const ApiLoginProvider int64 = 1

type UserApi struct {
	rp       repo.Repository
	validate *validator.Validate
	tokenMan *secure.TokenManager
}

func (a *UserApi) New(rpo repo.Repository, t *secure.TokenManager) {
	a.rp = rpo
	a.validate = validator.New()
	a.tokenMan = t
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

		// Get the user
		resp, err := a.rp.GetUserById(id)
		if err != nil {
			return c.JSON(http.StatusNotFound, &er.ErrResponse{er.ErrContent{er.ErrorUserNotFound, err.Error()}})
		}
		// Get the user profile
		resp.Profile, err = a.rp.GetUserProfileByUserId(resp.ID)
		if err != nil {
			return c.JSON(http.StatusNotFound, &er.ErrResponse{er.ErrContent{er.ErrorUserProfileNotFound, err.Error()}})
		}
		// Get the user security
		resp.Security, err = a.rp.GetSecurityInfoByUserId(resp.ID)
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

		ur := &models.User{Name: u.Name, Email: u.Email, Active: true, Security: &models.UserSecurity{LastMachine: c.RealIP()}}

		err = a.rp.InsertUser(ur)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, &er.ErrResponse{er.ErrContent{http.StatusInternalServerError, err.Error()}})
		}

		// Get all user information
		ur, err = a.rp.GetUserById(ur.ID)
		if err != nil {
			return c.JSON(http.StatusNotFound, &er.ErrResponse{er.ErrContent{er.ErrorUserNotFound, err.Error()}})
		}
		// Get the user security
		ur.Security, err = a.rp.GetSecurityInfoByUserId(ur.ID)
		if err != nil {
			return c.JSON(http.StatusNotFound, &er.ErrResponse{er.ErrContent{er.ErrorUserProfileNotFound, err.Error()}})
		}

		return c.JSON(http.StatusOK, ur)
	}
}

// Handler to POST User
// TODO: Get the user id from the token and it should match the same from the json object
func (a *UserApi) PostUser() echo.HandlerFunc {
	return func(c echo.Context) error {

		u := new(models.User)
		if err := c.Bind(u); err != nil {
			return c.JSON(http.StatusBadRequest, &er.ErrResponse{er.ErrContent{http.StatusBadRequest, err.Error()}})
		}
		if err := a.validate.Struct(u); err != nil {
			return c.JSON(http.StatusBadRequest, &er.ErrResponse{er.ErrContent{http.StatusBadRequest, err.Error()}})
		}

		if err := a.rp.UpdateUser(u); err != nil {
			return c.JSON(http.StatusBadRequest, &er.ErrResponse{er.ErrContent{http.StatusBadRequest, err.Error()}})
		}

		return c.JSON(http.StatusOK, u)
	}
}

// Handler to Login User
func (a *UserApi) LoginUser() echo.HandlerFunc {
	return func(c echo.Context) error {

		u := new(models.LoginUser)
		if err := c.Bind(u); err != nil {
			return c.JSON(http.StatusBadRequest, &er.ErrResponse{er.ErrContent{http.StatusBadRequest, err.Error()}})
		}
		if err := a.validate.Struct(u); err != nil {
			return c.JSON(http.StatusBadRequest, &er.ErrResponse{er.ErrContent{http.StatusBadRequest, err.Error()}})
		}

		// Get the user
		resp, err := a.rp.GetUserByEmail(u.Email)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, &er.ErrResponse{er.ErrContent{http.StatusUnauthorized, "Invalid credentials"}})
		}
		// Get the user profile
		resp.Profile, _ = a.rp.GetUserProfileByUserId(resp.ID)

		// Get the user security
		resp.Security, err = a.rp.GetSecurityInfoByUserId(resp.ID)
		if err != nil {
			return c.JSON(http.StatusNotFound, &er.ErrResponse{er.ErrContent{er.ErrorUserProfileNotFound, err.Error()}})
		}
		//Set the last machine
		resp.Security.LastMachine = c.RealIP()

		// Validate user password
		if !a.checkPasswordHash(u.Password, resp.Security.Hash) {
			return c.JSON(http.StatusUnauthorized, &er.ErrResponse{er.ErrContent{http.StatusUnauthorized, "Invalid credentials password"}})
		}

		// Create the JWT Token
		tc := &tok.TokenClaims{Email: resp.Email, ID: resp.ID}
		token, err := a.tokenMan.CreateToken(tc, "")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, &er.ErrResponse{er.ErrContent{er.ErrorCreatingToken, err.Error()}})
		}

		//Set the token in the Header
		c.Response().Header().Set(echo.HeaderAuthorization, token)

		//Update the LastMachine and LastLogin in a new go routine
		go func() {
			if err := a.rp.UpdateLoginData(resp.Security); err != nil {
				c.Logger().Errorf(err.Error())
			}
		}()

		return c.JSON(http.StatusOK, resp)
	}
}

// hashPassword Generates the hash of a given user password
func (a *UserApi) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// checkPasswordHash Validates that the user password is correct
func (a *UserApi) checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
