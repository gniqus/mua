package mua

import (
	"net/http"
)

type Context struct {
	Writer     http.ResponseWriter
	Request    *http.Request
	Path       string
	Method     string
	StatusCode int
	Params	   map[string]string
}

func (c *Context) FormValue(key string) string {
	return c.Request.FormValue(key)
}

func (c *Context) UrlValue(key string) string {
	return c.Request.URL.Query().Get(key)
}
