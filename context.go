package bfwork

import "net/http"

// Context 上下文
type Context struct {
	// 响应
	W http.ResponseWriter
	// 请求
	R *http.Request
	// Method
	Method string
	// 请求url
	Pattern string
	// 参数路由参数
	params map[string]string
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		W:       w,
		R:       r,
		Method:  r.Method,
		Pattern: r.URL.Path,
	}
}
