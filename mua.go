package mua

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
)

type Handler func(ctx *Context)

type engine struct {
	ctrl   string
	groups map[string]*routeGroup

	*router
}

var eng *engine
var once sync.Once

func GetEngine() *engine {
	once.Do(func() {
		eng = &engine{
			router: newRouter(),
			ctrl:   "/",
			groups: make(map[string]*routeGroup),
		}
		eng.groups[eng.ctrl] = &routeGroup{
			prefix:      eng.ctrl,
			middlewares: make([]Handler, 0),
		}
	})
	return eng
}

func (e *engine) Group(prefix string) *engine {
	e.ctrl = prefix
	if e.groups[e.ctrl] == nil {
		e.groups[e.ctrl] = &routeGroup{
			prefix:      e.ctrl,
			middlewares: make([]Handler, 0),
		}
	}
	return eng
}

func (e *engine) Use(middlewares ...Handler) *engine {
	e.groups[e.ctrl].middlewares = append(e.groups[e.ctrl].middlewares, middlewares...)
	return eng
}

func (e *engine) GET(path string, handler Handler) {
	e.registerRoute("GET", e.ctrl+path, handler)
}

func (e *engine) POST(path string, handler Handler) {
	e.registerRoute("POST", e.ctrl+path, handler)
}

func (e *engine) Start(address string) error {
	return http.ListenAndServe(address, e)
}

func (e *engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	method := strings.ToUpper(r.Method)
	node, params := e.findRoute(method, r.URL.Path)
	if node != nil {
		path := "/" + strings.Join(node.path, "/")
		middlewares := make([]Handler, 0)
		for prefix, group := range e.groups {
			if strings.HasPrefix(r.URL.Path, prefix) {
				middlewares = append(middlewares, group.middlewares...)
			}
		}
		middlewares = append(middlewares, e.handlers[method][path])
		newContext(w, r, params, middlewares).Next()
	} else {
		fmt.Fprintln(w, "404 NOT FOUND")
	}
}
