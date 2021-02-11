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
	if e.groups[eng.ctrl] == nil {
		e.groups[eng.ctrl] = &routeGroup{
			prefix:      eng.ctrl,
			middlewares: make([]Handler, 0),
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
			Writer:  w,
			Request: r,
			Path:    r.URL.Path,
			Method:  method,
			Params:  params,
		})
	} else {
		fmt.Fprintln(w, "404 NOT FOUND")
	}
}
