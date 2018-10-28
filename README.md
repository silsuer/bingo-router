# bingo-router

基于httprouter的路由模块

## 简介

这个包是为了自己的go web项目 [Bingo](https://github.com/silsuer/bingo) 所开发的模块之一，负责核心路由功能，这里将其抽出来作为单独包，
可以轻松的在其他项目中使用，而无需依赖 [Bingo](https://github.com/silsuer/bingo) 框架


特性:

 - [x] 路由快速查找

 - [x] 动态路由传参

 - [x] 子路由支持(路由组)

 - [x] 中间件支持

`bingo-router` 基于 [httprouter](https://github.com/julienschmidt/httprouter) , 性能强劲， `httprouter` 拥有的特性，`bingo-router`也有，并且参考 [Laravel](https://github.com/laravel/laravel) 的 [Pipeline](https://laravel-china.org/articles/2769/laravel-pipeline-realization-of-the-principle-of-single-component) 功能，为 `bingo-router` 添加中间件支持，可以快速实现拦截器等常见功能

后续准备加入的功能

 -  会话管理

 -  路由跳转

 -  服务的平滑重启和关闭

 -  工具方法补充(e.g.封装json返回等方法)

 -  测试代码补充

## 安装

```
go get -u github.com/silsuer/bingo-router
```

> bingo-router 使用 [glide](https://github.com/Masterminds/glide) 管理依赖，如果遇到无法 go install 问题，请先安装 `glide` 后，去 `$PATH/src/github.com/silsuer/bingo-router` 目录下执行 `glide install` 命令安装依赖


## 使用

`bingo-router` 已内置在 `[Bingo](https://github.com/silsuer/bingo)` 框架中，开箱即用，对于第三方项目:

1. 初始化路由器对象 `router`:

 ```go
   	// 创建一个路由器
   	router := bingo_router.New()
 ```

 路由器 `Router` 上挂载路由 `Route`，一个服务只能存在一个 `Router`，可存在多个`Route`，不要弄混

2. 创建路由

 ```go
   route := bingo_router.NewRoute().Get("/").Target(func(c *bingo_router.Context) {
   		fmt.Fprint(c.Writer,"hello! bingo-router!")
   	})
 ```

 目前上下文对象`context`只是对 `httprouter` 中对象的一个简单封装,用法与 `httprouter` 完全一致:

 ```go
   type Context struct {
   	Writer  http.ResponseWriter // 响应
   	Request *http.Request       // 请求
   	Params  Params              //参数
   }
 ```

 `NewRoute` 方法将会返回一个路由指针，可以调用:

 - `Get()` 方法，设定为Get方法

 - `Post()` 方法，设定为Get方法

 - `Patch()` 方法，设定为Get方法

 - `Put()` 方法，设定为Get方法

 - `Delete()` 方法，设定为Get方法

 - `Handle()` 方法，设定为Get方法

 - `Target()` 方法，指定该路由所对应的方法

 - `Middleware()` 方法，设定一个将在其中设置的中间件

 - `MiddlewareGroup()` 方法，批量设置中间件

 - `Prefix()` 方法，设置路由前缀，该前缀将对该路由的所有子路由有效，对当前路由无效

 - `Mount()` 方法，设置子路由

3. 将路由挂载到路由器中

 ```go
   router.Mount(route)
 ```

 路由器 `Router` 的

4. 开启服务器

 ```go
   http.ListenAndServe(":8080", r)
 ```

5. 在浏览器中访问 `http://localhost:8080`, 可以看到输出 `hello, bingo-router!`

## 进阶

### 1. 使用中间件

  使用 `NewRoute()` 方法之后将会获得一个 `Route` 对象指针，采用责任链模式，随意进行链式操作:

  ```go

	r2 := bingo_router.NewRoute().Get("/test").Middleware(func(c *bingo_router.Context, next func(c *bingo_router.Context)) {
		fmt.Fprint(c.Writer,"middleware")
		next(c)
	}).Target(func(c *bingo_router.Context) {
		fmt.Fprintln(c.Writer,"target")
	})

  ```

  中间件是一种 `MiddlewareHandle` 类型的方法 `func(c *Context, next func(c *Context))`

  `Middleware` 方法传入一个回调函数，该函数接受两个参数，第一个参数`context` 是上下文，是封装的`http.ResponseWriter`,`*http.Request`和`httprouter.Params`

  第二个参数是 `next` 方法，用来指定何时跳入下一个中间件，主要是用来实现前置中间件和后置中间件功能，当执行 `next` 方法后，将会跳入下一个中间件，直到所有中间件运行完毕后，才会继续执行下面的代码

  例如，上面的`r2`的效果是输出:

  ```
    middleware
    target
  ```

  如果将中间件中的`next()` 方法向上提一行的话:

  ```go
      r2 := bingo_router.NewRoute().Get("/test").Middleware(func(c *bingo_router.Context, next func(c *bingo_router.Context)) {
            next(c)
            fmt.Fprint(c.Writer,"middleware")
        }).Target(func(c *bingo_router.Context) {
            fmt.Fprintln(c.Writer,"target")
        })
  ```
  这样将会输出:

  ```go
     target
     middleware
  ```

  这样就可以实现后置中间件的效果了，可以做一些访问完控制器方法之后的后续工作

  当然，如果不在中间件中调用 `next()` 方法，请求就将在这个中间件中终止。


### 2. 使用中间件组

  `Route` 对象的 `MiddlewareGroup()` 方法可以实现批量注册中间件的功能，传入一个 `MiddlwareHandle` 的数组即可，不再赘述

### 3. 使用子路由和路由前缀

  在 `Laravel` 中，我们经常会用到路由组的概念，可以方便的进行整体控制，在 `bingo-router` 中，我也实现了这样的功能，但是`bingo-router`中没有路由组的概念，可以使用子路由实现类似效果,

  调用 `Route` 对象的 `Mount` 方法添加子路由:

  ```go

    bingo_router.NewRoute().Prefix("/api").Middleware(MiddlewareFunction).Mount(func(b *bingo_router.Builder) {
    		b.NewRoute().Get("/mount").Target(func(c *bingo_router.Context) {
    			fmt.Fprint(c.Writer, "hello mount")
    		})

    		b.NewRoute().Get("/mount2").Target(func(c *bingo_router.Context) {
                			fmt.Fprint(c.Writer, "hello mount2")
            })
    	})
  ```
  `Prefix` 方法，传入一个字符串，这个字符串将作用在所有子路由路径前，例如上面的代码，我们最终的访问路径就是 `/api/mount`和`/api/mount2`

  `Mount` 方法传入一个 `Builder` 指针，我们可以调用 `Builder`对象的 `NewRoute()` 方法为当前路由添加子路由,例如上面的代码，我们在访问 `/api/mount`和`/api/mount2`时都会经过我们指定的`MiddlewareFunction` 中间件


> 这个包还没有完全开发完成，后续会继续补充，也会在这两天内内置到 [Bingo](https://github.com/silsuer/bingo) 框架中

> 请大家多多关注主框架 [Bingo](https://github.com/silsuer/bingo) , 欢迎 STAR 、 PR


