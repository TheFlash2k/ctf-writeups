package fileserver

import (
	"html/template"
	"io"
	"log"
	"net/http"

	"workspace/fileserver/handlers"
	"workspace/fileserver/middlewares"
	"workspace/utils"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func StartFileHostingServer() {

	e := echo.New()
	t := &Template{}

	db := utils.CreateDB("fileserver.db")
	e.Renderer = t
	e.HTTPErrorHandler = customHTTPErrorHandler

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "FILESERVER: method=${method}, uri=${uri}, status=${status}\n",
	}))
	e.Use(middleware.Recover())
	e.Use(middlewares.ContextDB(db))

	e.Use(session.Middleware(sessions.NewCookieStore([]byte("[REDACTED]"))))

	e.Use(middlewares.CheckSession)
	e.Static("/static", "static")

	e.GET("/files", handlers.ListFiles)
	e.GET("/file", handlers.ViewFile)
	e.GET("/files/search", handlers.SearchFilesHandler)

	e.POST("/file/delete", handlers.DeleteFile)
	e.POST("/file", handlers.FileUploadHandler)

	e.GET("/file/remote", handlers.DownloadView)
	e.POST("/file/remote", handlers.DownloadHandler)

	e.GET("/login", handlers.LoginView)
	e.POST("/login", handlers.LoginHandler)

	e.POST("/profile/update", handlers.UpdateProfile)
	e.GET("/profile", handlers.ViewProfile)

	e.POST("/plugins/activate", handlers.ActivatePlugin)
	e.POST("/plugins/run", handlers.RunPlugin)
	e.GET("/plugins/all", handlers.AllPlugins)
	e.GET("/plugins", handlers.UserPlugins)
	e.HTTPErrorHandler = customHTTPErrorHandler

	// Start the file hosting server
	if err := e.Start(":8083"); err != nil {
		log.Println(err)
	}
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if t.templates == nil {
		t.templates = template.Must(template.ParseGlob("templates/*.html"))
	}

	return t.templates.ExecuteTemplate(w, name, data)
}

func customHTTPErrorHandler(err error, c echo.Context) {
	he, ok := err.(*echo.HTTPError)

	if ok {
		if he.Code == http.StatusNotFound {
			c.Redirect(http.StatusSeeOther, "/login")
			return
		}
	}

	c.Echo().DefaultHTTPErrorHandler(err, c)
}
