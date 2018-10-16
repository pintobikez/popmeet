package main

import (
	"context"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/color"
	"github.com/labstack/gommon/log"
	"github.com/pintobikez/popmeet/api"
	uti "github.com/pintobikez/popmeet/config"
	cnfs "github.com/pintobikez/popmeet/config/structures"
	er "github.com/pintobikez/popmeet/errors"
	rep "github.com/pintobikez/popmeet/repository"
	mysql "github.com/pintobikez/popmeet/repository/mysql"
	"gopkg.in/urfave/cli.v1"
	"io"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var (
	repo rep.Repository
)

const (
	StatusAvailable   = "Available"
	StatusUnavailable = "Unavailable"
)

// Start Http Server
func Handler(c *cli.Context) error {

	//
	apiInterest := new(api.InterestApi)

	// Echo instance
	e := echo.New()
	e.HTTPErrorHandler = serverErrorHandler
	e.Logger.SetLevel(log.INFO)
	e.Logger.SetOutput(LoadFileWriter(c.String("log-folder") + "/app.log"))

	e.Use(mw.Recover())
	e.Use(mw.Secure())
	e.Use(mw.RequestID())
	e.Pre(mw.RemoveTrailingSlash())

	//loads db connection
	dbConfig := new(cnfs.DatabaseConfig)
	err := uti.LoadConfigFile(c.String("database-file"), dbConfig)
	if err != nil {
		e.Logger.Fatal(err)
	}

	repo, err = mysql.New(dbConfig)
	if err != nil {
		e.Logger.Fatal(err)
	}

	// Database connect
	err = repo.Connect()
	if err != nil {
		e.Logger.Fatal(err)
	}
	defer repo.Disconnect()

	corsGET := mw.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.OPTIONS, echo.HEAD},
	}

	// Routes => healh
	e.GET("/health", healthStatus(repo), mw.CORSWithConfig(corsGET))

	// Routes => interests api
	apiInterest.SetRepository(repo)
	e.GET("/interest/all", apiInterest.GetAllInterest(), mw.CORSWithConfig(corsGET))
	e.GET("/interest/:id", apiInterest.GetInterest(), mw.CORSWithConfig(corsGET))

	// Start server
	colorer := color.New()
	colorer.Printf("⇛ %s service - %s\n", appName, color.Green(version))

	go func() {
		if err := start(e, c); err != nil {
			colorer.Printf(color.Red("⇛ shutting down the server\n"))
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	return nil
}

// Start http server
func start(e *echo.Echo, c *cli.Context) error {

	if c.String("ssl-cert") != "" && c.String("ssl-key") != "" {
		return e.StartTLS(
			c.String("listen"),
			c.String("ssl-cert"),
			c.String("ssl-key"),
		)
	}

	return e.Start(c.String("listen"))
}

// ServerErrorHandler sets the format of the error to be return by the server
func serverErrorHandler(err error, c echo.Context) {

	code := http.StatusServiceUnavailable
	msg := http.StatusText(code)

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		msg = he.Message.(string)
	}

	if c.Echo().Debug {
		msg = err.Error()
	}

	content := map[string]interface{}{
		"id":      c.Response().Header().Get(echo.HeaderXRequestID),
		"message": msg,
		"status":  code,
	}

	c.Logger().Errorj(content)

	if !c.Response().Committed {
		if c.Request().Method == echo.HEAD {
			c.NoContent(code)
		} else {
			c.JSON(code, &er.ErrResponse{er.ErrContent{code, msg}})
		}
	}
}

// File Retrieve the io.Writer from a file if exists, otherwise returns a os.Stdout
func LoadFileWriter(filePath string) io.Writer {
	file, err := os.OpenFile(
		filePath,
		os.O_RDWR|os.O_CREATE|os.O_APPEND,
		0666,
	)

	if err != nil {
		return os.Stdout
	}

	return file
}

// Handler for Health Status
func healthStatus(rp rep.Repository) echo.HandlerFunc {
	return func(c echo.Context) error {

		resp := new(er.HealthStatus)

		if err := rp.Health(); err != nil {
			resp.Repo.Status = StatusUnavailable
			resp.Repo.Detail = err.Error()
		}

		return c.JSON(http.StatusOK, resp)
	}
}
