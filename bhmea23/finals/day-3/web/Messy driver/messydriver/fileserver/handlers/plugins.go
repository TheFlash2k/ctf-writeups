package handlers

import (
	"bytes"
	"database/sql"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"unicode"
	"workspace/fileserver/dev"
	"workspace/fileserver/middlewares"
	"workspace/fileserver/models"

	"github.com/labstack/echo/v4"
	"github.com/mrz1836/go-sanitize"
)

var (
	httpInfoFeatureNotImplemented = "This feature is not implemented yet"
	httpInfoPluginSuccess         = "File converted succesfully"
	httpErrPluginExecFailed       = "Error running plugin"
)

func ActivatePlugin(c echo.Context) error {
	return c.String(http.StatusOK, httpInfoFeatureNotImplemented)
}

func RunPlugin(c echo.Context) error {

	username := c.Get("username").(string)
	user_directory := c.Get("user_directory").(string)

	pluginName := sanitize.PathName(c.FormValue("plugin_name"))

	filename := sanitize.PathName(c.FormValue("filename"))

	base_dir := os.Getenv("BASE_DIR")
	filepath := path.Join(filepath.Join(base_dir, "uploads"), user_directory, filename)

	db := c.Get("db").(*middlewares.DB)

	if !checkIfUserHasPluginActivated(db.Db, username, pluginName) {
		return c.String(http.StatusForbidden, "You cannot run this plugin")
	}

	err := dev.RunPlugin(path.Join("plugins/", titleToSnakeCase(pluginName)), filepath)

	if err != nil {
		return c.String(http.StatusInternalServerError, httpErrPluginExecFailed)
	}

	return c.String(http.StatusOK, httpInfoPluginSuccess)
}

func AddPlugin(c echo.Context) error {

	return c.String(http.StatusOK, httpInfoFeatureNotImplemented)
}

func AllPlugins(c echo.Context) error {
	base_dir := os.Getenv("BASE_DIR")
	files, err := os.ReadDir(filepath.Join(base_dir, "plugins/"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to read plugins directory")
	}

	plugins := make([]string, 0)
	for _, f := range files {
		// Exclude directories and only consider .so files
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".so") {
			pluginName := strings.TrimSuffix(f.Name(), ".so")
			formattedName := snakeCaseToTitle(pluginName)
			plugins = append(plugins, formattedName)
		}
	}

	return c.JSON(http.StatusOK, plugins)
}

func UserPlugins(c echo.Context) error {
	username := c.Get("username").(string) // assuming the username is set in the context

	db := c.Get("db").(*middlewares.DB)

	// Retrieve the activated_plugins string from the database for the specific user
	var activatedPlugins string
	err := db.Db.QueryRow("SELECT activated_plugin FROM users WHERE username = ?", username).Scan(&activatedPlugins)
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, "User not found")
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to retrieve plugins from DB")
	}

	// Split the string into individual plugin names
	plugins := strings.Split(activatedPlugins, ",")

	// Return the list of plugins as JSON
	return c.JSON(http.StatusOK, plugins)
}

func checkIfUserHasPluginActivated(db *sql.DB, username string, pluginName string) bool {
	user, err := models.GetUserByUsername(db, username)
	if err != nil {

		// Handle error (maybe log it or return false)
		return false
	}

	for _, plugin := range user.ActivatedPlugin {
		if plugin == pluginName {
			return true
		}
	}

	return false
}

func snakeCaseToTitle(s string) string {
	parts := strings.Split(s, "_")
	for i := range parts {
		parts[i] = strings.Title(parts[i])
	}
	return strings.Join(parts, " ")
}
func titleToSnakeCase(str string) string {
	var buffer bytes.Buffer
	words := splitWords(str)

	for i, word := range words {
		if i > 0 {
			buffer.WriteString("_")
		}
		buffer.WriteString(strings.ToLower(word))
	}

	return buffer.String()
}

// splitWords splits a string into words separated by spaces or camel case.
func splitWords(str string) []string {
	var words []string

	wordStart := -1
	for i, ch := range str {
		if unicode.IsUpper(ch) {
			if wordStart != -1 {
				words = append(words, str[wordStart:i])
			}
			wordStart = i
		} else if ch == ' ' || ch == '_' {
			if wordStart != -1 {
				words = append(words, str[wordStart:i])
				wordStart = -1
			}
		}
	}

	if wordStart != -1 {
		words = append(words, str[wordStart:])
	}

	return words
}
