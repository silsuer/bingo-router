package bingo_router

import (
	"net/http"
	"strings"
	"strconv"
	"github.com/modood/table"
)

type Handle func(http.ResponseWriter, *http.Request, Params)

type Param struct {
	Key   string
	Value string
}

type Params []Param

func (ps Params) ByName(name string) string {
	for i := range ps {
		if ps[i].Key == name {
			return ps[i].Value
		}
	}
	return ""
}

type Router struct {
	trees map[string]*node

	RedirectTrailingSlash bool

	RedirectFixedPath bool

	HandleMethodNotAllowed bool

	HandleOPTIONS bool

	NotFound http.HandlerFunc

	MethodNotAllowed http.Handler

	PanicHandler func(http.ResponseWriter, *http.Request, interface{})
}

var _ http.Handler = New()

func New() *Router {
	return &Router{
		RedirectTrailingSlash:  true,
		RedirectFixedPath:      true,
		HandleMethodNotAllowed: true,
		HandleOPTIONS:          true,
	}
}

// GET is a shortcut for router.Handle("GET", path, handle)
func (r *Router) GET(path string, route *Route) {
	r.Handle("GET", path, route)
}

// HEAD is a shortcut for router.Handle("HEAD", path, handle)
func (r *Router) HEAD(path string, route *Route) {
	r.Handle("HEAD", path, route)
}

// OPTIONS is a shortcut for router.Handle("OPTIONS", path, handle)
func (r *Router) OPTIONS(path string, route *Route) {
	r.Handle("OPTIONS", path, route)
}

// POST is a shortcut for router.Handle("POST", path, handle)
func (r *Router) POST(path string, route *Route) {
	r.Handle("POST", path, route)
}

// PUT is a shortcut for router.Handle("PUT", path, handle)
func (r *Router) PUT(path string, route *Route) {
	r.Handle("PUT", path, route)
}

// PATCH is a shortcut for router.Handle("PATCH", path, handle)
func (r *Router) PATCH(path string, route *Route) {
	r.Handle("PATCH", path, route)
}

// DELETE is a shortcut for router.Handle("DELETE", path, handle)
func (r *Router) DELETE(path string, route *Route) {
	r.Handle("DELETE", path, route)
}

// Handle registers a new request handle with the given path and method.
//
// For GET, POST, PUT, PATCH and DELETE requests the respective shortcut
// functions can be used.
//
// This function is intended for bulk loading and to allow the usage of less
// frequently used, non-standardized or custom methods (e.g. for internal
// communication with a proxy).
func (r *Router) Handle(method, path string, route *Route) {
	if path[0] != '/' {
		panic("path must begin with '/' in path '" + path + "'")
	}

	if r.trees == nil {
		r.trees = make(map[string]*node)
	}

	root := r.trees[method]
	if root == nil {
		root = new(node)
		r.trees[method] = root
	}

	// 将指针转换成对象
	root.addRoute(path, *route)
}

// HandlerFunc is an adapter which allows the usage of an http.HandlerFunc as a
// request handle.
func (r *Router) HandlerFunc(method, path string, handler http.HandlerFunc) {
	r.Handler(method, path, handler)
}

func (r *Router) Handler(method, path string, handler http.Handler) {
	route := &Route{}
	route.targetMethod = func(context *Context) {
		handler.ServeHTTP(context.Writer, context.Request)
	}
	r.Handle(method, path, route)
}

// ServeFiles serves files from the given file system root.
// The path must end with "/*filepath", files are then served from the local
// path /defined/root/dir/*filepath.
// For example if root is "/etc" and *filepath is "passwd", the local file
// "/etc/passwd" would be served.
// Internally a http.FileServer is used, therefore http.NotFound is used instead
// of the Router's NotFound handler.
// To use the operating system's file system implementation,
// use http.Dir:
//     router.ServeFiles("/src/*filepath", http.Dir("/var/www"))
func (r *Router) ServeFiles(path string, root http.FileSystem) {
	if len(path) < 10 || path[len(path)-10:] != "/*filepath" {
		panic("path must end with /*filepath in path '" + path + "'")
	}

	fileServer := http.FileServer(root)
	route := &Route{}
	route.targetMethod = func(c *Context) {
		c.Request.URL.Path = c.Params.ByName("filepath")
		fileServer.ServeHTTP(c.Writer, c.Request)
	}
	r.GET(path, route)
}

func (r *Router) recv(w http.ResponseWriter, req *http.Request) {
	if rcv := recover(); rcv != nil {
		r.PanicHandler(w, req, rcv)
	}
}

// Lookup allows the manual lookup of a method + path combo.
// This is e.g. useful to build a framework around this router.
// If the path was found, it returns the handle function and the path parameter
// values. Otherwise the third return value indicates whether a redirection to
// the same path with an extra / without the trailing slash should be performed.
func (r *Router) Lookup(method, path string) (Route, Params, bool) {
	if root := r.trees[method]; root != nil {
		return root.getValue(path)
	}
	return Route{}, nil, false
}

func (r *Router) allowed(path, reqMethod string) (allow string) {
	if path == "*" { // server-wide
		for method := range r.trees {
			if method == "OPTIONS" {
				continue
			}

			// add request method to list of allowed methods
			if len(allow) == 0 {
				allow = method
			} else {
				allow += ", " + method
			}
		}
	} else { // specific path
		for method := range r.trees {
			// Skip the requested method - we already tried this one
			if method == reqMethod || method == "OPTIONS" {
				continue
			}

			route, _, _ := r.trees[method].getValue(path)
			if route.targetMethod != nil {
				// add request method to list of allowed methods
				if len(allow) == 0 {
					allow = method
				} else {
					allow += ", " + method
				}
			}
		}
	}
	if len(allow) > 0 {
		allow += ", OPTIONS"
	}
	return
}

//func (r *Router) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
//	fmt.Fprint(ctx, "hello")
//}

// ServeHTTP makes the router implement the http.Handler interface.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if r.PanicHandler != nil {
		defer r.recv(w, req)
	}

	path := req.URL.Path

	if root := r.trees[req.Method]; root != nil {

		path := req.URL.Path

		route, ps, tsr := root.getValue(path)

		if route.targetMethod != nil {

			// 这里可以封装上下文
			context := &Context{w, req, ps}

			// 建立管道，执行中间件最终到达路由
			new(Pipeline).Send(context).Through(route.middleware).Then(func(context *Context) {
				route.targetMethod(context)
			})

			return
		} else if req.Method != "CONNECT" && path != "/" {
			code := 301 // Permanent redirect, request with GET method
			if req.Method != "GET" {
				// Temporary redirect, request with same method
				// As of Go 1.3, Go does not support status code 308.
				code = 307
			}

			if tsr && r.RedirectTrailingSlash {
				if len(path) > 1 && path[len(path)-1] == '/' {
					req.URL.Path = path[:len(path)-1]
				} else {
					req.URL.Path = path + "/"
				}
				http.Redirect(w, req, req.URL.String(), code)
				return
			}

			// Try to fix the request path
			if r.RedirectFixedPath {
				fixedPath, found := root.findCaseInsensitivePath(
					CleanPath(path),
					r.RedirectTrailingSlash,
				)
				if found {
					req.URL.Path = string(fixedPath)
					http.Redirect(w, req, req.URL.String(), code)
					return
				}
			}
		}
	}

	if req.Method == "OPTIONS" && r.HandleOPTIONS {
		// Handle OPTIONS requests
		if allow := r.allowed(path, req.Method); len(allow) > 0 {
			w.Header().Set("Allow", allow)
			return
		}
	} else {
		// Handle 405
		if r.HandleMethodNotAllowed {
			if allow := r.allowed(path, req.Method); len(allow) > 0 {
				w.Header().Set("Allow", allow)
				if r.MethodNotAllowed != nil {
					r.MethodNotAllowed.ServeHTTP(w, req)
				} else {
					http.Error(w,
						http.StatusText(http.StatusMethodNotAllowed),
						http.StatusMethodNotAllowed,
					)
				}
				return
			}
		}
	}

	// Handle 404
	if r.NotFound != nil {
		r.NotFound.ServeHTTP(w, req)
	} else {
		http.NotFound(w, req)
	}
}

//var path []string
var prefix []string
var middlewares map[string][]MiddlewareHandle
var currentPointer int // 当前是第几层路由

func (r *Router) Mount(routes ...*Route) {
	prefix = []string{}
	//currentPointer = 0
	middlewares = make(map[string][]MiddlewareHandle)
	for _, route := range routes {
		r.MountRoute(route)
	}

}

// 向其中挂载路由
func (r *Router) MountRoute(route *Route) {

	// 拼接路径
	// 挂载中间件
	// middlewares 变量保存第x层路由的中间件
	// 这些中间件对当前路由和当前路由的子路由将产生作用
	setMiddlewares(currentPointer, route)

	p := getPrefix(currentPointer) + route.path // 当前路径是所有前缀数组连接在一起，加上当前路由的path
	// 如果一个路由设置了前缀，则这个前缀会作用在所有的子路由上
	prefix = append(prefix, route.prefix)

	// 设置中间件
	//setRouteMiddlewares(currentPointer, route)
	if route.method != "" && p != "" {
		r.Handle(route.method, p, route)
	}

	// 如果路由有子路由，则将子路由挂载进去，如果没有，
	if len(route.mount) > 0 {
		for _, subRoute := range route.mount {
			currentPointer += 1 // 添加一层
			r.MountRoute(subRoute)
		}
	} else {
		if currentPointer > 0 {
			currentPointer -= 1 // 减小一层
		}
		// 如果没有子路由，则清空掉该临时路径和中间件数组
		//path = ""
		//middlewares = []MiddlewareHandle{}
	}

}

// 根据当前是第几层路由，获取前缀
func getPrefix(current int) string {
	if len(prefix) > current-1 && len(prefix) != 0 {
		return strings.Join(prefix[:current], "")
	}
	return ""
}

// 设置中间件，根据当前是第x层路由，将前面的路由放入当前路由中
func setMiddlewares(current int, route *Route) {
	key := "p" + strconv.Itoa(currentPointer)
	for _, v := range route.middleware {
		middlewares[key] = append(middlewares[key], v)
	}

	// 将当前路由的父路由的都放入当前路由中
	for i := 0; i < currentPointer; i++ {
		key = "p" + strconv.Itoa(i)
		if list, ok := middlewares[key]; ok {
			for _, v := range list {
				route.middleware = append(route.middleware, v)
			}
		}
	}
}

type Output struct {
	Method     string
	URI        string
	Name       string
	Action     string
	Middleware string
}

// 在控制台中打印所有路由
func (r *Router) PrintRoutes() {
	var outputs []Output
	// 遍历所有路由，拼接成table结构体，然后打印输出
	for m := range r.trees {
		rr := r.trees[m].route
		if rr.path == "" {
			continue
		}

		// 调用路由的一个方法，传入一个 output数组，如果有
		outputs = append(outputs, rr.print(outputs)...)
	}

	//fmt.Println(outputs)
	table.Output(outputs)
}
