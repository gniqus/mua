package mua

import (
	"strings"
)

type node struct {
	sons   []*node
	blurry bool
	value  string
	isTail bool
	path   []string
}

type trie struct {
	root *node
}

func newTrie() *trie {
	return &trie{
		root: &node{
			sons: make([]*node, 0),
			path: make([]string, 0),
			value: "/",
			blurry: false,
			isTail: false,
		},
	}
}

func (t *trie) insert(path string) {
	t.insertHelp(t.root, t.split(path), 0)
}

func (t *trie) search(path string) *node {
	return t.searchHelp(t.root, t.split(path), 0)
}

func (t *trie) split(str string) []string {
	values := strings.Split(str, "/")
	result := make([]string, 0)
	for _, value := range values {
		if value == "" {
			continue
		}
		result = append(result, value)
	}
	return result
}

func (t *trie) find(node *node, value string) *node {
	for _, node := range node.sons {
		if node.value == value {
			return node
		}
	}
	return nil
}

func (t *trie) isBlurry(value string) bool {
	return len(value) > 0 && (value[0] == ':' || value[0] == '*')
}

func (t *trie) insertHelp(cur *node, values []string, level int) {
	if len(values) == level || strings.HasPrefix(cur.value, "*") {
		cur.isTail = true
		return
	}
	value := values[level]
	next := t.find(cur, value)
	if next == nil {
		next =  &node{
			sons:   make([]*node, 0),
			blurry: t.isBlurry(value),
			value:  value,
			path:   append(cur.path, value),
		}
		cur.sons = append(cur.sons, next)
	}
	t.insertHelp(next, values, level+1)
}

func (t *trie) searchHelp(cur *node, values []string, level int) (result *node) {
	if len(values) == level || strings.HasPrefix(cur.value, "*") {
		if cur.isTail {
			result = cur
		}
	} else {
		value := values[level]
		for i := 0; i < len(cur.sons) && result == nil; i++ {
			if cur.sons[i].value == value || cur.sons[i].blurry {
				result = t.searchHelp(cur.sons[i], values, level+1)
			}
		}
	}
	return result
}

//func (t *trie) searchHelp(cur *node, values []string, level int) *node {
//	if len(values) == level || strings.HasPrefix(cur.value, "*") {
//		if cur.isTail {
//			return cur
//		}
//		return nil
//	}
//	value := values[level]
//	for _, son := range cur.sons {
//		if son.value == value || son.blurry {
//			result := t.searchHelp(son, values, level+1)
//			if result != nil {
//				return result
//			}
//		}
//	}
//	return nil
//}
