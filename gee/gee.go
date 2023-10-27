package gee

import (
	"log"
	"net/http"
	"strings"
)

// HandlerFunc 是一个类型，它是一个函数类型，实际上是一个适配器（Adapter）
type HandlerFunc func(*Context)

// Engine 引擎实现ServeHTTP接口
type Engine struct {
	*RouterGroup
	router *router
	groups []*RouterGroup
}

type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc //支持中间件
	parent      *RouterGroup  //支持嵌套
	engine      *Engine       //所有组共享一个引擎实例
}

// New 创造一个gee
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine

}

// Group 组被定义为创建新的RouterGroup
// 请记住，所有组共享同一个Engine实例
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

// 添加接口(HandleFunc 是一个函数，它用于注册一个请求处理函数（handler Function）)
func (group *RouterGroup) addRouter(method string, comp string, handler HandlerFunc) {
	path := group.prefix + comp
	log.Printf("Route %4s - %s", method, path)
	group.engine.router.addRouter(method, path, handler)

}

// 添加中间件
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

func (group *RouterGroup) POST(path string, handler HandlerFunc) {
	group.addRouter("POST", path, handler)
}

func (group *RouterGroup) DELETE(path string, handler HandlerFunc) {
	group.addRouter("DELETE", path, handler)
}

func (group *RouterGroup) PUT(path string, handler HandlerFunc) {
	group.addRouter("PUT", path, handler)
}

func (group *RouterGroup) GET(path string, handler HandlerFunc) {
	group.addRouter("GET", path, handler)
}

// Run 定义启动http服务器的方法
func (engine *Engine) Run(address string) (err error) {
	return http.ListenAndServe(address, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, req)
	c.handlers = middlewares
	engine.router.handle(c)
}
