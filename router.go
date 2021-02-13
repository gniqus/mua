package mua

import (
	"errors"
	"strings"
)

type router struct {
	tries    map[string]*trie
	handlers map[string]map[string]Handler
}

type routeGroup struct {
	prefix      string
	parent      *routeGroup
	middlewares []Handler
}

func newRouter() *router {
	router := &router{
		tries:    make(map[string]*trie),
		handlers: make(map[string]map[string]Handler),
	}
	router.tries["POST"] = newTrie()
	router.tries["GET"] = newTrie()
	router.handlers["POST"] = make(map[string]Handler)
	router.handlers["GET"] = make(map[string]Handler)
	return router
}

func (r *router) registerRoute(method string, path string, handler Handler) error {
	if r.handlers[method] == nil || r.tries[method] == nil {
		return errors.New("method error")
	}
	if len(path) > 1 && path[0] == '/' && path[1] == '/' {
		path = path[1:]
	}
	r.tries[method].insert(path)
	r.handlers[method][path] = handler
	return nil
}

func (r *router) findRoute(method string, path string) (*node, map[string]string) {
	result := r.tries[method].search(path)
	params := make(map[string]string)
	r.fillParams(result, r.tries[method].split(path), params)
	return result, params
}

func (r *router) fillParams(node *node, values []string, params map[string]string) {
	if node == nil {
		return
	}
	for i, value := range node.path {
		if len(value) > 1 && strings.HasPrefix(value, ":") {
			params[value[1:]] = values[i]
		}
		if len(value) > 1 && strings.HasPrefix(value, "*") {
			params[value[1:]] = strings.Join(values[i:], "/")
		}
	}
}
