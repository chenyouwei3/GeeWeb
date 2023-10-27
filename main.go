package main

import (
	"gee"
	"net/http"
)

func main() {
	r := gee.New()
	r.GET("/index", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>cyw</h1>")
	})
	c1 := r.Group("/c")
	{
		c1.GET("/c1", func(c *gee.Context) {
			c.HTML(http.StatusOK, "<h1>hello</h1>")
		})

		c1.GET("/c2", func(c *gee.Context) {

		})
	}

	r.Run(":9999")
}
