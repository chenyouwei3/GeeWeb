package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	Writer     http.ResponseWriter
	Req        *http.Request
	Path       string
	Method     string
	Params     map[string]string
	StatusCode int
	handlers   []HandlerFunc
	index      int
	engine     *Engine
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		index:  -1,
	}
}

//func (c *Context) Next() {
//	c.index++            // 将index加1
//	s := len(c.handlers) // 获取handlers的长度并赋值给s
//	for ; c.index < s; c.index++ {
//		c.handlers[c.index](c) // 调用handlers[index]并将c作为参数传入
//	}
//}

func (c *Context) Next() {
	for len(c.handlers) != 0 {
		handlerFunc := c.handlers[0] // 获取handlers的第一个元素
		c.handlers = c.handlers[1:]  // 移除handlers的第一个元素
		handlerFunc(c)               // 调用handlerFunc并将c作为参数传入
	}
}

func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
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

// 方法用于返回纯文本格式的响应
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

// Data 方法用于返回原始字节数据格式的响应。它设置状态码，并将传入的字节数据写入到响应中
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

// HTML 方法用于返回 HTML 格式的响应。它设置响应的 Content-Type 为 text/html，
// 设置状态码，并将传入的 HTML 字符串转换为字节写入到响应中
func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}
