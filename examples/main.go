package main

import (
	"fmt"
	"github.com/silsuer/bingo-router"
	"net/http"
)

// 定义
func setRoute() *bingo_router.Route {
	return bingo_router.NewRoute().Prefix("/common").Middleware(func(c *bingo_router.Context, next func(c *bingo_router.Context)) {
		fmt.Fprintln(c.Writer, "common middleware 1")
		next(c)
	}).Middleware(func(c *bingo_router.Context, next func(c *bingo_router.Context)) {
		fmt.Fprintln(c.Writer, "common middleware 2")
		next(c)
	}).Mount(func(b *bingo_router.Builder) {
		b.NewRoute().Get("/test1").Target(func(c *bingo_router.Context) {
			fmt.Fprint(c.Writer, "hello test 1")
		}).Middleware(func(c *bingo_router.Context, next func(c *bingo_router.Context)) {
			next(c)
			fmt.Fprintln(c.Writer, " middleware test 1")

		})

		b.NewRoute().Get("/test2").Target(func(c *bingo_router.Context) {
			fmt.Fprint(c.Writer, "hello test2")
		}).Middleware(func(c *bingo_router.Context, next func(c *bingo_router.Context)) {
			next(c)
			fmt.Fprintln(c.Writer, "test 2 middleware 1")

		}).Middleware(func(c *bingo_router.Context, next func(c *bingo_router.Context)) {
			fmt.Fprintln(c.Writer, "test 2 middleware 2")
			next(c)
		})
	})
}

func main() {
	// 创建一个路由
	r := bingo_router.New()

	// 得到一个路由
	rr := setRoute()

	//r2 := bingo_router.NewRoute().Get("/test").Middleware(func(c *bingo_router.Context, next func(c *bingo_router.Context)) {
	//	fmt.Fprint(c.Writer,"middleware")
	//	next(c)
	//}).Target(func(c *bingo_router.Context) {
	//	fmt.Fprintln(c.Writer,"target")
	//})
	// 挂载路由方法，遍历传入的所有路由，然后遍历每个路由里的子路由，并设定路径等，设定中间件等
	r.Mount(rr)

	// 启动服务
	http.ListenAndServe(":8080", r)

	// TODO 服务的平滑重启和关闭
}
