package mua

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request
	Path    string
	Method  string
	Params  map[string]string
}

func (c *Context) FormValue(key string) string {
	return c.Request.FormValue(key)
}

func (c *Context) UrlValue(key string) string {
	return c.Request.URL.Query().Get(key)
}

func (c *Context) EchoString(format string, values ...interface{}) {
	fmt.Fprintf(c.Writer, format, values...)
}

func (c *Context) EchoJSON(obj interface{}) {
	data, _ := json.Marshal(obj)
	c.Writer.Write(data)
}

func (c *Context) EchoData(data []byte) {
	c.Writer.Write(data)
}
