package mua

import (
	"fmt"
	"net/http"
)

type Handler func(http.ResponseWriter, *http.Request)

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
	if handler, ok := e.handlers[r.Method][r.URL.Path]; ok {
		handler(w, r)
	} else {
		fmt.Fprintln(w, "404 NOT FOUND")
	}
}
