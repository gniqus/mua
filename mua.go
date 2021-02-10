package mua

import (
	"fmt"
	"net/http"
	"strings"
)

type Handler func(ctx *Context)

type engine struct {
	*router
}

var eng *engine

func GetEngine() *engine {
	if eng == nil {
		eng = &engine{
			router: newRouter(),
		}
	}
	return eng
}

func (e *engine) GET(path string, handler Handler) {
	e.registerRoute("GET", path, handler)
}

func (e *engine) POST(path string, handler Handler) {
	e.registerRoute("POST", path, handler)
}

func (e *engine) Start(address string) error {
	return http.ListenAndServe(address, e)
}

func (e *engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	method := strings.ToUpper(r.Method)
	node, params := e.findRoute(method, r.URL.Path)
	if node != nil {
		path := "/" + strings.Join(node.path, "/")
		e.handlers[method][path](&Context{
			Writer: w,
			Request: r,
			Path: r.URL.Path,
			Method: method,
			Params: params,
		})
	} else {
		fmt.Fprintln(w, "404 NOT FOUND")
	}
}
