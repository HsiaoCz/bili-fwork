package bfwork

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// HandleFunc 视图函数签名
// type HandleFunc func(w http.ResponseWriter, r *http.Request)
type HandleFunc func(c *Context)

// 一个server需要什么功能
// 1.启动
// 2.关闭
// 3.注册路由的方法

// 为什么要抽象这个server呢?
// 有些网站走http协议，有些网站走https协议
// 为了兼容不同的协议
type server interface {
	// 硬性要求，必须组合http.Handler
	http.Handler
	// 启动服务
	Start(addr string) error
	// 关闭服务
	Stop() error
	// 注册路由的方法
	// 非常核心的方法
	// 这个方法不能被外界使用
	addRouter(method string, pattern string, handleFunc HandleFunc)
}

// 选项模式

type HTTPOption func(h *HTTPServer)

type HTTPServer struct {
	srv  *http.Server
	stop func() error
	// routers 临时存放的路由的位置
	routers map[string]HandleFunc
}

// 路由的设计
// "GET-login":HandleFunc1,
// "POST-login":HandleFunc2,

func WithHTTPServerStop(fn func() error) HTTPOption {
	return func(h *HTTPServer) {
		if fn == nil {
			fn = func() error {
				fmt.Println("111111111")
				// os.Signal类型的channel
				quit := make(chan os.Signal)
				// 如果匹配到中断信息，会将信号传到channel里面
				// 如果没有匹配到，channel会一直阻塞
				signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
				<-quit
				log.Println("shutdown Server....")

				// 创建一个超时的上下文，在这里等五秒钟
				// 等待任务执行完毕
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				if err := h.srv.Shutdown(ctx); err != nil {
					log.Fatal("server Shutdown", err)
				}
				// 关闭之后执行的操作
				select {
				case <-ctx.Done():
					log.Println("timeout of 5 seconds...")
				}
				return nil
			}
		}
		h.stop = fn
	}
}

func NewHTTP(opts ...HTTPOption) *HTTPServer {
	h := &HTTPServer{
		routers: map[string]HandleFunc{},
	}
	for _, opt := range opts {
		opt(h)
	}
	return h
}

// 接收请求转发请求
// ServeHTTP方法向前对接前端请求，向后对接咱们的框架
func (h *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 1.匹配路由
	key := fmt.Sprintf("%s-%s", r.Method, r.URL.Path)
	handler, ok := h.routers[key]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 NOT FOUND"))
		return
	}
	// 2.构造当前请求的上下文
	c := NewContext(w, r)
	log.Printf("request %s-%s", c.Method, c.Pattern)
	// 2.转发请求
	handler(c)
}

// GET方法
func (h *HTTPServer) GET(pattern string, handleFunc HandleFunc) {
	h.addRouter(http.MethodGet, pattern, handleFunc)
}

// POST请求
func (h *HTTPServer) POST(pattern string, handleFunc HandleFunc) {
	h.addRouter(http.MethodPost, pattern, handleFunc)
}

// PUT请求
func (h *HTTPServer) PUT(pattern string, handleFunc HandleFunc) {
	h.addRouter(http.MethodPut, pattern, handleFunc)
}

// Delete请求
func (h *HTTPServer) DELETE(pattern string, handleFunc HandleFunc) {
	h.addRouter(http.MethodDelete, pattern, handleFunc)
}

func (h *HTTPServer) Start(addr string) error {
	h.srv = &http.Server{
		Addr:    addr,
		Handler: h,
	}
	log.Println("the server is running on port", addr)
	return h.srv.ListenAndServe()
}

// 这里服务的关闭需要优雅的关闭
func (h *HTTPServer) Stop() error {
	return h.stop()
}

// 实现路由注册的方法
// 注册路由的时机，项目启动的时候，启动之后就不能注册了
// 问题，注册路由，注册到哪里
func (h *HTTPServer) addRouter(method string, pattern string, handleFunc HandleFunc) {
	// 这里构建唯一的key
	key := fmt.Sprintf("%s-%s", method, pattern)
	h.routers[key] = handleFunc
}
