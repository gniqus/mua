package mua

import (
	"errors"
	"strings"
)

type router struct {
	trie     map[string]*trie
	handlers map[string]map[string]Handler
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
	method = strings.ToUpper(method)
	if r.handlers[method] == nil || r.trie[method] == nil {
		return errors.New("method error")
	}
	r.trie[method].Insert(path)
	r.handlers[method][path] = handler
	return nil
}

func (r *router) findRoute(method string, path string) (*node, map[string]string) {
	values := strings.Split(path, "/")
	node := r.trie[method].Search(path)
	params := make(map[string]string)
	if node != nil {
		for i, value := range node.path {
			if strings.HasPrefix(value, ":") {
				params[value[1:]] = values[i]
			}
			if strings.HasPrefix(value, "*") {
				params[value[1:]] = strings.Join(values[i:], "/")
				break
			}
		}
	}
	return node, params
}
