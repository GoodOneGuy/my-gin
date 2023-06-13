package my_gin

import (
	"log"
	"net/http"
	"testing"
	"time"
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
	r.GET("/index", func(c *Context) {
		c.HTML(http.StatusOK, "<h1>Index Page</h1>")
	})
	v1 := r.Group("/v1")
	{
		v1.GET("/", func(c *Context) {
			c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
		})

		v1.GET("/hello", func(c *Context) {
			// expect /hello?name=geektutu
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})
	}
	v2 := r.Group("/v2")
	{
		v2.GET("/hello/:name", func(c *Context) {
			// expect /hello/geektutu
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
		v2.POST("/login", func(c *Context) {
			c.JSON(http.StatusOK, H{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})

	}

	r.Run(":9999")
}

func onlyForV2() []HandlerFunc {
	return []HandlerFunc{
		func(c *Context) {
			// Start timer
			c.Next()
			t := time.Now()
			// Calculate resolution time
			log.Printf("[%d] %s in %v middleware1 for group v2\n", c.StatusCode, c.Req.RequestURI, time.Since(t))
		},
		func(c *Context) {
			// Start timer
			t := time.Now()
			// Calculate resolution time
			log.Printf("[%d] %s in %v middleware2 for group v2\n", c.StatusCode, c.Req.RequestURI, time.Since(t))
			c.Next()
		},
	}
}

func TestEngine_MyGin(t *testing.T) {
	r := New()
	r.Use(Logger()) // global midlleware
	r.GET("/", func(c *Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})

	v2 := r.Group("/v2")
	v2.Use(onlyForV2()...) // v2 group middleware
	{
		v2.GET("/hello/:name", func(c *Context) {
			// expect /hello/geektutu
			log.Printf("real handler\n")
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
	}

	r.Run(":9999")
}
