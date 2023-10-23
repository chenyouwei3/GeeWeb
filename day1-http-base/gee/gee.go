package gee

import (
	"fmt"
	"net/http"
)

// HandlerFunc 是一个类型，它是一个函数类型，实际上是一个适配器（Adapter）
type HandlerFunc func(http.ResponseWriter, *http.Request)

// Engine 引擎实现ServeHTTP接口
type Engine struct {
	router map[string]HandlerFunc
}

// New 创造一个gee
func New() *Engine {
	return &Engine{router: make(map[string]HandlerFunc)}
}

// 添加接口(HandleFunc 是一个函数，它用于注册一个请求处理函数（Handler Function）)
func (engine *Engine) addRouter(method string, path string, handler HandlerFunc) {
	key := method + "-" + path
	engine.router[key] = handler
}

func (engine *Engine) POST(path string, handler HandlerFunc) {
	engine.addRouter("POST", path, handler)
}

func (engine *Engine) DELETE(path string, handler HandlerFunc) {
	engine.addRouter("DELETE", path, handler)
}

func (engine *Engine) PUT(path string, handler HandlerFunc) {
	engine.addRouter("PUT", path, handler)
}

func (engine *Engine) GET(path string, handler HandlerFunc) {
	engine.addRouter("GET", path, handler)
}

// Run 定义启动http服务器的方法
func (engine *Engine) Run(address string) (err error) {
	return http.ListenAndServe(address, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := req.Method + "-" + req.URL.Path
	if handler, ok := engine.router[key]; ok {
		handler(w, req)
	} else {
		fmt.Fprintf(w, "404 NOT FOUND:%s\n", req.URL)
	}
}
