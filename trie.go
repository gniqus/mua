package mua

import "strings"

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
		},
	}
}

func (t *trie) Insert(path string) {
	values := strings.Split(path, "/")
	t.insert(t.root, values[1:], 0)
}

func (t *trie) Search(path string) *node {
	values := strings.Split(path, "/")
	return t.search(t.root, values[1:], 0)
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

func (t *trie) insert(cur *node, values []string, level int) {
	if len(values) == level {
		return
	}
	value := values[level]
	next := t.find(cur, value)
	if next == nil {
		cur.sons = append(cur.sons, &node{
			path:   append(cur.path, value),
			sons:   make([]*node, 0),
			blurry: t.isBlurry(value),
			value:  value,
			isTail: len(values)-1 == level+1,
		})
	}
	t.insert(next, values, level+1)
}

func (t *trie) search(cur *node, values []string, level int) *node {
	if len(values) == level || strings.HasPrefix(cur.value, "*") {
		if cur.isTail {
			return cur
		}
		return nil
	}

	value := values[level]
	for _, son := range cur.sons {
		if son.value != value {
			continue
		}
		result := t.search(son, values, level+1)
		if result != nil {
			return result
		}
	}
	return nil
}
