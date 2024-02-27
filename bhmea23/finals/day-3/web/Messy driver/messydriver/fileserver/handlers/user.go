package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"workspace/fileserver/middlewares"
	"workspace/fileserver/models"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

var (
	httpErrUserNotFound          = "Error finding user, what did you do?? Please report the bug to me"
	httpErrRequiredFieldsMissing = "Required fields are missing"
)

type UserLogin struct {
	Password string `json:"password"`
	ClientId string `json:"clientId"`
}

func LoginView(c echo.Context) error {

	return c.Render(http.StatusOK, "login.html", nil)
}

func LoginHandler(c echo.Context) error {

	username := c.FormValue("username")
	password := c.FormValue("password")

	if username == "" || password == "" {
		return c.NoContent(http.StatusUnauthorized)
	}

	jsonBody, err := json.Marshal(&middlewares.PasswordBody{
		Password: password,
	})

	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	url := url.URL{
		Scheme:   "http",
		Host:     "localhost:8080",
		Path:     "/",
		RawQuery: "host=localhost&port=8082&path=/auth" + url.QueryEscape("?SECRET="+os.Getenv("SECRET")),
	}

	resp, err := http.Post(url.String(), "application/json", bytes.NewBuffer(jsonBody))

	if err != nil {
		return c.Redirect(http.StatusUnauthorized, "/login")
	}

	defer resp.Body.Close()

	// Check response status

	if resp.StatusCode == http.StatusOK {

		Db := c.Get("db").(*middlewares.DB)
		user, err := models.GetUserByUsername(Db.Db, username)

		if err != nil {
			return c.String(http.StatusInternalServerError, httpErrUserNotFound)

		}

		sess, err := session.Get("session", c)

		sess.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   86400 * 7,
			HttpOnly: true,
		}

		sess.Values["user_directory"] = user.Directory
		sess.Values["logged_in"] = true
		sess.Values["username"] = user.Username
		err = sess.Save(c.Request(), c.Response())
		if err != nil {
			return c.String(http.StatusInternalServerError, "Re-login")
		}

		if err != nil {
			return c.String(http.StatusInternalServerError, httpErrUserNotFound)
		}

		return c.Redirect(http.StatusFound, "/files")
	} else {

		return c.Redirect(http.StatusFound, "/login")
	}

}

type Profile struct {
	Nickname    string `form:"nickname"`
	Nationality string `form:"nationality"`
}

func UpdateProfile(c echo.Context) error {
	db := c.Get("db").(*middlewares.DB)

	var newProfile Profile
	if err := c.Bind(&newProfile); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid input")
	}

	if newProfile.Nickname == "" || newProfile.Nationality == "" {
		return c.JSON(http.StatusBadRequest, httpErrRequiredFieldsMissing)
	}

	username := c.Get("username").(string)

	query := fmt.Sprintf(
		"UPDATE profiles SET nickname = %s, nationality = '%s' WHERE id = %d",
		sanitizeSQL(newProfile.Nickname),
		sanitizeSQL(newProfile.Nationality),
		1)

	_, err := db.Db.Exec(query, newProfile.Nickname, newProfile.Nationality, username)

	if err != nil {

		return c.JSON(http.StatusInternalServerError, "Failed to update profile")
	}

	return c.JSON(http.StatusOK, "Profile updated successfully")
}

func ViewProfile(c echo.Context) error {
	user_directory := c.Get("user_directory").(string)
	db := c.Get("db").(middlewares.DB)

	var profile Profile
	err := db.Db.QueryRow("SELECT nickname, nationality FROM profiles WHERE id = ?", user_directory).Scan(&profile.Nickname, &profile.Nationality)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to retrieve profile")
	}

	userDirectory := filepath.Join("uploads", user_directory)
	files, err := os.ReadDir(userDirectory)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to read user files directory")
	}
	fileCount := len(files) // Get the number of files

	data := struct {
		Profile   Profile
		FileCount int
	}{
		Profile:   profile,
		FileCount: fileCount,
	}

	return c.Render(http.StatusFound, "profile.html", data)
}

func sanitizeSQL(sqlx string) string {
	blacklist := []string{
		"DROP", "SELECT", "INSERT", "DELETE", "--", "/*", "*/", "EXEC", "EXECUTE",
	}

	for _, keyword := range blacklist {
		if strings.Contains(strings.ToUpper(sqlx), keyword) {
			sqlx = strings.ReplaceAll(sqlx, keyword, "")
		}
	}

	return sqlx
}
