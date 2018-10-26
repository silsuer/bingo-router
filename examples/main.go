package main
import (
	"github.com/silsuer/bingo-router"
	"fmt"
	"net/http"
)

// 定义
func setRoute() *bingo_router.Route {
	return bingo_router.NewRoute().Prefix("aaa").Mount(func(b *bingo_router.Builder) {
		b.NewRoute().Get("/").Target(func(c *bingo_router.Context) {
			fmt.Fprint(c.Writer, "hello")
		})
	})
}

func main() {
	// 为了避免协程太多，我们大部分使用连接池
	// 创建一个路由
	// 给路由配置选项
	// 配置连接池
	//r := httprouter.New()
	//r := bingo_router.New()
	r := bingo_router.New()

	// 先改成传入Route对象，就可以
	rr := setRoute()
	//fmt.Println(rr)
	// 挂载路由方法，遍历传入的所有路由，然后遍历每个路由里的子路由，并设定路径等，设定中间件等
	r.Mount(rr)

	//r.GET("/", func(writer http.ResponseWriter, request *http.Request, params bingo_router.Params) {
	//	fmt.Fprint(writer, "hello111")
	//})
	//
	http.ListenAndServe(":8080", r)

	//fasthttp.ListenAndServe(":8080", r)
	//fasthttp.Serve(":8080",)
}
