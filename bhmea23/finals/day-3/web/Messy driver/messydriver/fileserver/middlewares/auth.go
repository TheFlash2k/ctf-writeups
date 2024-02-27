package middlewares

import (
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type PasswordBody struct {
	Password string `json:"password"`
}

func CheckSession(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {

		if c.Path() == "/login" || strings.HasPrefix(c.Path(), "/static") || strings.HasPrefix(c.Path(), "/file/remote") {
			return next(c)
		}

		sess, err := session.Get("session", c)

		if err != nil {
			return c.Redirect(http.StatusFound, "/login")
		}

		val, isok := sess.Values["logged_in"].(bool)

		if isok && val {

			val, _ := sess.Values["user_directory"].(string)
			c.Set("user_directory", val)
			val, _ = sess.Values["username"].(string)
			c.Set("username", val)

			return next(c)
		}

		err = sess.Save(c.Request(), c.Response())
		if err != nil {
			log.Fatal("failed to delete session", err)
		}

		return c.Redirect(http.StatusFound, "/login")
	}
}
