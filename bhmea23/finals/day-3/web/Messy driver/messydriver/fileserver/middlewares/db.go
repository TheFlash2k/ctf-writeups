package middlewares

import (
	"database/sql"

	"github.com/labstack/echo/v4"
)

type DB struct {
	Db *sql.DB
}

func ContextDB(db *sql.DB) echo.MiddlewareFunc {

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("db", &DB{db})
			return next(c)
		}
	}
}
