package bingo_router

// 控制器接口
type IController interface {
	Index(c *Context)   // 对应 get : path方法
	Create(c *Context)  // 对应 get : path/create 方法
	Store(c *Context)    // 对应 post : path 方法
	Update(c *Context)  // 对应put/patch : path/:id   方法
	Edit(c *Context)    // 对应get : path/:id/edit
	Show(c *Context)    // 对应 get: route/:id
	Destroy(c *Context) // 对应 delete: route/:id
}

// 默认基本控制器
type Controller struct {
}

func (cc *Controller) Index(c *Context) {

}

func (cc *Controller) Create(c *Context) {

}

func (cc *Controller) Store(c *Context) {

}

func (cc *Controller) Update(c *Context) {

}

func (cc *Controller) Edit(c *Context) {

}
func (cc *Controller) Show(c *Context) {

}

func (cc *Controller) Destroy(c *Context) {

}
