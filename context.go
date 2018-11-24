package bingo_router

import (
	"github.com/json-iterator/go"
	"mime/multipart"
	"net/http"
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

// 获取get表单的参数
func (c *Context) Get(key string) string {
	return c.Request.FormValue(key)
}

// 获取post表单的参数
func (c *Context) Post(key string) string {
	return c.Request.PostFormValue(key)
}

func (c *Context) File(key string) (multipart.File, *multipart.FileHeader, error) {
	return c.Request.FormFile(key)
}

func (c *Context) GetWithDefault(key, def string) string {
	if c.Request.FormValue(key) == "" {
		return def
	}
	return c.Request.FormValue(key)
}

func (c *Context) PostWithDefault(key, def string) string {
	if c.Request.PostFormValue(key) == "" {
		return def
	}
	return c.Request.PostFormValue(key)
}
