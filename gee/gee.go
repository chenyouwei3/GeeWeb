package gee

import (
	"net/http"
)

// HandlerFunc 是一个类型，它是一个函数类型，实际上是一个适配器（Adapter）
type HandlerFunc func(*Context)

// Engine 引擎实现ServeHTTP接口
type Engine struct {
	router *router
}

// New 创造一个gee
func New() *Engine {
	return &Engine{router: newRouter()}
}

// 添加接口(HandleFunc 是一个函数，它用于注册一个请求处理函数（handler Function）)
func (engine *Engine) addRouter(method string, path string, handler HandlerFunc) {
	engine.router.addRouter(method, path, handler)
}

func (engine *Engine) POST(path string, handler HandlerFunc) {
	engine.addRouter("POST", path, handler)
}

func (engine Engine) DELETE(path string, handler HandlerFunc) {
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
	c := newContext(w, req)
	engine.router.handle(c)
}
