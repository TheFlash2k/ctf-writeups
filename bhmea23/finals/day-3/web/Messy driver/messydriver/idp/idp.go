package idp

import (
	"log"
	"net/http"
	"os"
	"workspace/utils"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type User struct {
	password string
	username string
}

type AuthUser struct {
	Password string `query:"password" json:"password"`
	ClientId string `json:"clientId" header:"Client-Id" query:"clientId" `
}

type IDPApiSecret struct {
	Secret string `json:"SECRET"`
}

func StartIDPServer() {

	e := echo.New()
	// Define the routes and handlers for the IDP server

	db := utils.CreateDB("idp.db")

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "IDP: method=${method}, uri=${uri}, status=${status}\n",
	}))

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			secret := c.QueryParam("SECRET")

			if secret == os.Getenv("SECRET") {
				return next(c)
			}

			return c.NoContent(http.StatusUnauthorized)
		}
	})

	e.Any("/auth", func(c echo.Context) error {

		var authUser AuthUser

		err := c.Bind(&authUser)

		if err != nil || authUser.Password == "" {
			return c.NoContent(http.StatusUnauthorized)
		}

		var user User

		err = db.QueryRow("SELECT * from users where password = '%s'", authUser.Password).Scan(&user)
		if err != nil {
			return c.NoContent(http.StatusUnauthorized)
		}

		return c.String(200, "User is correct he can go in")
	})
	e.HTTPErrorHandler = customHTTPErrorHandler

	// Start the IDP server
	if err := e.Start(":8082"); err != nil {
		log.Println(err)
	}
}

func customHTTPErrorHandler(err error, c echo.Context) {
	he, ok := err.(*echo.HTTPError)
	if ok {
		if he.Code == http.StatusNotFound {
			c.Redirect(http.StatusSeeOther, "/auth")
			return
		}
	}

	c.Echo().DefaultHTTPErrorHandler(err, c)
}
