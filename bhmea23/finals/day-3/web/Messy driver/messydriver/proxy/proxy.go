package proxy

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/patrickmn/go-cache"
	"golang.org/x/exp/slices"
)

type CustomBinder struct{}

func (cb *CustomBinder) Bind(i interface{}, c echo.Context) (err error) {
	db := new(echo.DefaultBinder)

	if err := db.BindPathParams(c, i); err != nil {
		return err
	}

	if err = db.BindQueryParams(c, i); err != nil {
		return err
	}

	if err = db.BindHeaders(c, i); err != nil {
		return err
	}
	return db.BindBody(c, i)
}

type HTTPRequest struct {
	Host string `header:"X-HOST" query:"host"`
	Port string `header:"X-PORT" query:"port"`
	Path string `header:"X-PATH" query:"path" `
}

type CachedResponse struct {
	ResponseBody []byte
	StatusCode   int
	Headers      map[string]string
}

var (
	METHODS_TO_CACHE = []string{"GET", "OPTIONS", "HEAD", "TRACE"}
)

func StartProxy() {
	e := echo.New()

	ccache := cache.New(5*time.Minute, 6*time.Minute)

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "PROXY: method=${method}, uri=${uri}, status=${status}\n",
	}))

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {

		return func(c echo.Context) error {

			host := c.QueryParam("host")
			port := c.QueryParam("port")
			path := c.QueryParam("path")

			if strings.HasPrefix(path, "/") {
				path = strings.TrimLeft(path, "/")
			}

			URL := fmt.Sprintf("http://%s:%s/%s", host, port, strings.Split(path, "?")[0])

			parsedURL, err := url.Parse(URL)

			if err != nil {
				return c.String(http.StatusBadRequest, "BAD URL")
			}

			cacheKey := ComputeSHA256(parsedURL.String())

			value, found := ccache.Get(cacheKey)

			if found {
				cachedValue, ok := value.(CachedResponse)
				if !ok {
					ccache.Delete(parsedURL.String())
				} else {

					for k, v := range cachedValue.Headers {
						c.Response().Header().Set(k, v)
					}
					c.Blob(cachedValue.StatusCode, cachedValue.Headers["Content-Type"], cachedValue.ResponseBody)
					return nil
				}

			}
			c.Set("cache_key", cacheKey)

			return next(c)
		}

	})

	e.Use(func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var httpRequest HTTPRequest

			var path string

			binder := CustomBinder{}

			err := binder.Bind(&httpRequest, c)

			if err != nil {
				return c.String(http.StatusBadRequest, "Proxy is angry, please give him the needed infos")
			}

			if strings.HasPrefix(httpRequest.Path, "/") {
				path = strings.TrimPrefix(httpRequest.Path, "/")
			}

			URL := fmt.Sprintf("http://%s:%s/%s", httpRequest.Host, httpRequest.Port, path)

			_, err = url.Parse(URL)

			if err != nil {
				return c.String(http.StatusBadRequest, "BAD URL")
			}

			var body []byte
			_, err = c.Request().Body.Read(body)

			request, err := http.NewRequest(c.Request().Method, URL, bytes.NewBuffer(body))

			if err != nil {
				return c.NoContent(400)
			}

			client := &http.Client{}
			response, error := client.Do(request)

			if error != nil {
				return c.String(http.StatusNotFound, "Invalid request properties, the server you want to reach might be not accessible")
			}

			responseBody, error := io.ReadAll(response.Body)

			if error != nil {
				return c.String(http.StatusInternalServerError, "Error reading response")
			}

			c.Blob(response.StatusCode, response.Header.Get("Content-Type"), responseBody)

			if slices.Contains(METHODS_TO_CACHE, c.Request().Method) {
				headers := make(map[string]string, 0)
				for k, v := range response.Header {
					headers[k] = v[0]
					c.Response().Header().Set(k, v[0])
				}
				cacheKey := c.Get("cache_key").(string)
				ccache.Set(cacheKey, CachedResponse{
					ResponseBody: responseBody,
					StatusCode:   response.StatusCode,
					Headers:      headers,
				}, time.Second*5)
			}

			return nil

		}
	})

	e.Start(":8080")

}

func ComputeSHA256(input string) string {
	hash := sha256.New()
	hash.Write([]byte(input))
	return hex.EncodeToString(hash.Sum(nil))
}
