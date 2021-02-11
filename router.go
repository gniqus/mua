package mua

import (
	"errors"
	"strings"
)

type router struct {
	trie     map[string]*trie
	handlers map[string]map[string]Handler
}

type routerGroup struct {
	parent      *routerGroup
	prefix      string
	middlewares []Handler
}

func newRouter() *router {
	router := &router{
		trie:     make(map[string]*trie),
		handlers: make(map[string]map[string]Handler),
	}
	router.trie["POST"] = newTrie()
	router.trie["GET"] = newTrie()
	router.handlers["POST"] = make(map[string]Handler)
	router.handlers["GET"] = make(map[string]Handler)
	return router
}

func (r *router) registerRoute(method string, path string, handler Handler) error {
	if r.handlers[method] == nil || r.trie[method] == nil {
		return errors.New("method error")
	}
	r.trie[method].insert(path)
	r.handlers[method][path] = handler
	return nil
}

func (r *router) findRoute(method string, path string) (*node, map[string]string) {
	result := r.trie[method].search(path)
	params := make(map[string]string)
	r.fillParams(result, r.trie[method].split(path), params)
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
