package gee

import (
	"fmt"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	Writer     http.ResponseWriter
	Req        *http.Request
	Path       string
	Method     string
	StatusCode int
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
	}
}

// PostForm 方法用于获取请求中的表单值，根据传入的 key 返回相应的表单值
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

// Query 方法用于获取请求 URL 中的查询参数值，根据传入的 key 返回相应的查询参数值。Query 方法用于获取请求 URL 中的查询参数值，根据传入的 key 返回相应的查询参数值。
func (c *Context) Query(key string) string {
	//通过调用 URL 方法获取请求的 URL 对象，并调用 Query 方法获取查询参数的集合
	//最后调用 Get 方法通过指定的键名 key 获取对应的值并返回
	return c.Req.URL.Query().Get(key)
}

// Status 方法用于设置响应的状态码，并通过 Writer.WriteHeader 方法将状态码写入响应头。
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// SetHeader 方法用于设置响应头的键值对
func (c *Context) SetHeader(key string, value string) {
	//通过调用 Header 方法获取响应头对象
	//并调用 Set 方法设置指定键名 key 对应的值为 value
	c.Writer.Header().Set(key, value)
}

func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}
