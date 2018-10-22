package main

import (
	"github.com/silsuer/bingo-router"
	"net/http"
	"fmt"
)

// 定义

func main() {
	// 为了避免协程太多，我们大部分使用连接池
	// 创建一个路由
	// 给路由配置选项
	// 配置连接池
	//r := httprouter.New()
	//r := bingo_router.New()
	r := bingo_router.New()

	r.GET("/", func(writer http.ResponseWriter, request *http.Request, params bingo_router.Params) {
		fmt.Fprint(writer, "hello111")
	})

	http.ListenAndServe(":8080", r)

	//fasthttp.ListenAndServe(":8080", r)
	//fasthttp.Serve(":8080",)
}
