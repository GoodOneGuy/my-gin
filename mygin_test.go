package my_gin

import (
	"net/http"
	"testing"
)

func newTestRouter() *router {
	r := newRouter()
	r.addRoute("GET", "/", nil)
	r.addRoute("GET", "/hello/:name", nil)
	r.addRoute("GET", "/hello/b/c", nil)
	r.addRoute("GET", "/hi/:name", nil)
	r.addRoute("GET", "/assets/*filepath", nil)
	return r
}

func TestEngine_ServeHttp(t *testing.T) {
	r := New()
	r.GET("/", func(c *Context) {
		c.HTML(http.StatusOK, "<h1>Hello my gin</h1>")
	})

	r.GET("/hello", func(c *Context) {
		// expect /hello?name=test
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	r.GET("/hello/:name", func(c *Context) {
		// expect /hello/test
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	})

	r.POST("/hello/:name", func(c *Context) {
		// expect /hello/test
		c.String(http.StatusOK, "hello %s, you're post at %s\n", c.Param("name"), c.Path)
	})

	r.GET("/hello/:name/:id", func(c *Context) {
		// expect /hello/geektutu
		c.String(http.StatusOK, "hello %s, id %s, you're at %s\n", c.Param("name"), c.Param("id"), c.Path)
	})

	r.GET("/assets/*filepath", func(c *Context) {
		c.Json(http.StatusOK, H{"filepath": c.Param("filepath")})
	})

	r.Run(":9999")
}
