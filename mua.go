package mua

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
)

type Handler func(ctx *Context)

type engine struct {
	ctrl   string
	groups map[string]*routeGroup

	*router
	tmpl *template.Template
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

func (e *engine) Mapping(virtual string, real string) *engine {
	handler := http.StripPrefix(path.Join(e.ctrl, virtual), http.FileServer(http.Dir(real)))
	e.GET(path.Join(virtual, "/*filepath"), func(c *Context) {
		filepath := c.Params["filepath"]
		if _, err := os.Open(path.Join(real, filepath)); err != nil {
			c.EchoString("404 NOT FOUND")
			return
		}
		handler.ServeHTTP(c.Writer, c.Request)
	})
	return eng
}

func (e *engine) LoadTmpls(dir string, funcMap template.FuncMap) *engine {
	e.tmpl = template.Must(template.New("mua").Funcs(funcMap).ParseGlob(dir))
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

func (e *engine) GET(p string, handler Handler) {
	e.registerRoute("GET", path.Join(e.ctrl, p), handler)
}

func (e *engine) POST(p string, handler Handler) {
	e.registerRoute("POST", path.Join(e.ctrl, p), handler)
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
		ctx := newContext(w, r)
		ctx.Params = params
		ctx.Middlewares = middlewares
		ctx.Tmpl = e.tmpl
		ctx.Next()
	} else {
		fmt.Fprintln(w, "404 NOT FOUND")
	}
}
