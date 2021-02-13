package mua

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Context struct {
	Writer      http.ResponseWriter
	Request     *http.Request
	Path        string
	Method      string
	Params      map[string]string
	Middlewares []Handler
	Index       int
}

func newContext(w http.ResponseWriter, r *http.Request, params map[string]string, middlewares []Handler) *Context {
	return &Context{
		Writer:      w,
		Request:     r,
		Path:        r.URL.Path,
		Method:      strings.ToUpper(r.Method),
		Params:      params,
		Middlewares: middlewares,
		Index:       -1,
	}
}

func (c *Context) FormValue(key string) string {
	return c.Request.FormValue(key)
}

func (c *Context) UrlValue(key string) string {
	return c.Request.URL.Query().Get(key)
}

func (c *Context) Next() {
	c.Index++
	c.Middlewares[c.Index](c)
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
