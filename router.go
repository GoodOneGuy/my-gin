package my_gin

import (
	"net/http"
	"strings"
)

type HandlerFunc func(c *Context)

type router struct {
	roots    map[string]*trieNode
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*trieNode),
		handlers: make(map[string]HandlerFunc),
	}
}

func parsePattern(pattern string) []string {
	arr := strings.Split(pattern, "/")
	var paths []string
	for _, item := range arr {
		if item != "" {
			paths = append(paths, item)
			if item[0] == '*' {
				break
			}
		}
	}

	return paths
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {

	paths := parsePattern(pattern)
	key := method + "-" + pattern

	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &trieNode{}
	}

	r.roots[method].insert(pattern, paths, 0)
	r.handlers[key] = handler
}

func (r *router) handle(c *Context) {

	root, ok := r.roots[c.Method]
	if !ok {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		return
	}

	params := make(map[string]string)
	searchPaths := parsePattern(c.Path)
	node := root.search(searchPaths, 0)

	if node != nil {
		paths := parsePattern(node.pattern)
		for i, path := range paths {
			if path[0] == ':' {
				params[path[1:]] = searchPaths[i]
			}
			if path[0] == '*' && len(path) > 1 {
				params[path[1:]] = strings.Join(searchPaths[i:], "/")
			}
		}
	}

	if node != nil {
		c.Params = params
		key := c.Method + "-" + node.pattern
		c.handlers = append(c.handlers, r.handlers[key])
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}

	c.Next()
}
