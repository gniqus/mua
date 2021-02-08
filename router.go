package mua

import (
	"errors"
	"strings"
)

type router struct {
	handlers map[string]map[string]Handler
}

func newRouter() *router {
	router := &router{
		handlers: make(map[string]map[string]Handler),
	}
	router.handlers["GET"] = make(map[string]Handler)
	router.handlers["POST"] = make(map[string]Handler)
	return router
}

func (r *router) registerRoute(method string, path string, handler Handler) error {
	method = strings.ToUpper(method)
	if r.handlers[method] == nil {
		return errors.New("method error")
	}
	r.handlers[method][path] = handler
	return nil
}
