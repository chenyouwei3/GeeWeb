package gee

import (
	"html/template"
	"log"
	"net/http"
	"path"
	"strings"
)

// HandlerFunc 是一个类型，它是一个函数类型，实际上是一个适配器（Adapter）
type HandlerFunc func(*Context)

// Engine 引擎实现ServeHTTP接口,统一处理，调用router的handle，从map中根据url查找handler
// 再次封装，由groupRouter来实现路由的功能
type Engine struct {
	*RouterGroup
	router        *router
	groups        []*RouterGroup
	htmlTemplates *template.Template
	funcMap       template.FuncMap
}

// New 创造一个gee
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{
		engine: engine,
	}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine

}

type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc //支持中间件
	parent      *RouterGroup  //支持嵌套
	engine      *Engine       //所有组共享一个引擎实例
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

// 创建静态处理程序
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(group.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		//检查文件是否存在和/或我们是否有访问权限
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

// Static 提供静态文件
func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPath := path.Join(relativePath, "/*filepath")
	//注册get文件
	group.GET(urlPath, handler)
}

// 添加接口(HandleFunc 是一个函数，它用于注册一个请求处理函数（handler Function）)
func (group *RouterGroup) addRouter(method string, comp string, handler HandlerFunc) {
	path := group.prefix + comp
	log.Printf("Route %4s - %s", method, path)
	group.engine.router.addRouter(method, path, handler)

}

// Use 添加中间件
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

// SetFuncMap 方法接收一个类型为 template.FuncMap 的参数 funcMap，它用于设置模板引擎的自定义函数映射。
// template.FuncMap 是一个 map[string]interface{} 类型，用于将函数名映射到对应的函数
func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

// LoadHTMLGlob 函数根据指定的模式加载所有匹配的 HTML 模板，并将结果存储在 engine.htmlTemplates 字段中
func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
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
	c.engine = engine
	engine.router.handle(c)
}
