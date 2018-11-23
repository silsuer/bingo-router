package bingo_router

import (
	"net/http"
	"github.com/json-iterator/go"
)

// 上下文结构体
type Context struct {
	Writer  http.ResponseWriter // 响应
	Request *http.Request       // 请求
	Params  Params              //参数
	//Session Session             // 保存session
}

// 重定向，传入路径和状态码
func (c *Context) Redirect(path string) {
	url := c.Request.URL.Host + path
	http.Redirect(c.Writer, c.Request, url, http.StatusFound)
}

// 输出字符串
func (c *Context) String(data string) {
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write([]byte(data))
}

// 输出json
func (c *Context) ResponseJson(data interface{}) {
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	d, _ := json.Marshal(data)
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write([]byte(d))
}
