package bingo_router

import (
	"strings"
	"runtime"
	"reflect"
)

const GET = "GET"
const POST = "POST"
const DELETE = "DELETE"
const PUT = "PUT"
const PATCH = "PATCH"

// 路由
type Route struct {
	path         string             // 路径
	targetMethod TargetHandle       // 要执行的方法
	method       string             // 访问类型 是get post 或者其他
	alias        string             // 路由的别名，并没有什么卵用的样子.......
	name         string             // 路由名
	mount        []*Route           // 子路由
	middleware   []MiddlewareHandle // 挂载的中间件
	prefix       string             // 路由前缀，该前缀仅对子路由有效
}

// 路由的构建器
type Builder struct {
	routes []*Route // 根据这个builder创建的
}

func (b *Builder) NewRoute() *Route {
	r := NewRoute()
	b.routes = append(b.routes, r)
	return r
}

type TargetHandle func(c *Context)
type MiddlewareHandle func(c *Context, next func(c *Context))

func NewRoute() *Route {
	return &Route{}
}

// 添加路由时需要，设置为Get方法
func (r *Route) Get(path string) *Route {
	//return r.Request(GET, path, target)
	r.path = path
	r.method = GET
	return r
}

// 添加路由时需要，设置为Post方法
func (r *Route) Post(path string) *Route {
	//return r.Request(POST, path, target)
	r.path = path
	r.method = POST
	return r
}

// 添加路由时需要，设置为put方法
func (r *Route) Put(path string) *Route {
	//return r.Request(PUT, path, target)
	r.path = path
	r.method = PUT
	return r
}

// 添加路由时需要，设置为patch方法
func (r *Route) Patch(path string) *Route {
	//return r.Request(PATCH, path, target)
	r.path = path
	r.method = PATCH
	return r
}

// 添加路由时需要，设置为delete方法
func (r *Route) Delete(path string) *Route {
	//return r.Request(DELETE, path, target)
	r.path = path
	r.method = DELETE
	return r
}

// 传入一个控制器，自动构建多个路由
func (r *Route) Resource(path string, controller IController) *Route {
	// 在当前路由下挂载子路由
	r.Mount(func(b *Builder) {
		b.NewRoute().Get(path).Target(controller.Index)
		b.NewRoute().Get(path + "/detail/:id").Target(controller.Show)
		b.NewRoute().Post(path).Target(controller.Store)
		b.NewRoute().Get(path + "/create").Target(controller.Create)
		b.NewRoute().Put(path + "/:id").Target(controller.Update)
		b.NewRoute().Patch(path + "/:id").Target(controller.Update)
		b.NewRoute().Get(path + "/edit/:id").Target(controller.Edit)
		b.NewRoute().Delete(path + "/:id").Target(controller.Destroy)
	})
	return r
}

// 这里传入一个回调
func (r *Route) Target(target TargetHandle) *Route {
	r.name = r.path + "." + r.method

	return r.Request(r.method, r.path, target)
}

func (r *Route) Request(method string, path string, target TargetHandle) *Route {
	r.method = method
	r.path = path
	r.targetMethod = target
	return r
}

// 路由前缀，该前缀仅会对子路由有效，对当前路由无效
func (r *Route) Prefix(prefix string) *Route {
	r.prefix = prefix
	return r
}

func (r *Route) Middleware(h MiddlewareHandle) *Route {
	r.middleware = append(r.middleware, h)
	return r
}

func (r *Route) MiddlewareGroup(hg []MiddlewareHandle) *Route {
	for _, h := range hg {
		r.Middleware(h)
	}
	return r
}

func (r *Route) Name(name string) *Route {
	if r.name != "" {
		return r
	}
	r.name = name
	return r
}

// 挂载子路由，这里只是将回调中的路由放入
func (r *Route) Mount(rr func(b *Builder)) *Route {
	builder := new(Builder)
	rr(builder)
	// 遍历这个路由下建立的所有子路由，将路由放入父路由上
	for _, route := range builder.routes {
		r.mount = append(r.mount, route)
	}
	return r
}

func (r *Route) print(arr []Output) []Output {
	o := Output{}
	o.Name = r.name
	o.URI = r.path
	o.Method = r.method
	o.Middleware = strings.Join(r.printMiddlewares(), "|")
	o.Action = runtime.FuncForPC(reflect.ValueOf(r.targetMethod).Pointer()).Name()
	arr = append(arr, o)
	if len(r.mount) > 0 {
		for rr := range r.mount {
			arr = r.mount[rr].print(arr)
		}
	}

	return arr
}

// 返回所有的中间件名称组合成的数组
func (r *Route) printMiddlewares() []string {
	var res []string
	for m := range r.middleware {
		res = append(res, runtime.FuncForPC(reflect.ValueOf(r.middleware[m]).Pointer()).Name())
	}
	return res
}

// 每个请求进来都要生成一个管道，根据管道执行中间件最后到达目的路由
