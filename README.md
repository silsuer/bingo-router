# bingo-router
基于httprouter的路由模块

// 重写httprouter种的router结构体，添加中间件支持

httprouter 原来是添加一个handle

改为添加一个 Route结构体

// 添加中间件组
// 添加路由组
路由组： 主路由+子路有
前缀   // 在注册的时候让前缀和路由放入其中
中间件 // 中间件中包括其他中间件

支持中间件（前置中间件，后置中间件）
支持中间件组

支持路由组

Router 是路由器
Route是路由，向Router中添加Route

NewRoute().Prefix("ddd").addMiddlewares().Son(function(builder){
   builder.get().target().name()
})

NewGroup().addMiddlewares(middleware1,middleware2).addRoutes(function(builder){
  builder.get().target().name()
  
})

使用NewRoute创建



