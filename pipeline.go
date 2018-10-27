package bingo_router

type Pipeline struct {
	send    *Context           // 穿过管道的上下文
	through []MiddlewareHandle // 中间件数组
	current int                // 当前执行到第几个中间件
}

// new(Pipeline).send(context).through(middleware).then(function(context){})

func (p *Pipeline) Send(context *Context) *Pipeline {
	p.send = context
	return p
}

func (p *Pipeline) Through(middlewares []MiddlewareHandle) *Pipeline {
	p.through = middlewares
	return p
}

func (p *Pipeline) Exec() {
	if len(p.through) > p.current {
		m := p.through[p.current]
		p.current += 1
		m(p.send, func(c *Context) {
			p.Exec()
		})
	}

}

// 这里是路由的最后一站
func (p *Pipeline) Then(then func(context *Context)) {
	// 按照顺序执行
	// 将then作为最后一站的中间件
	var m MiddlewareHandle
	m = func(c *Context, next func(c *Context)) {
		then(c)
		next(c)
	}
	p.through = append(p.through, m)
	p.Exec()
}
