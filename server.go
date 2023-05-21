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

// 一个server需要什么功能
// 1.启动
// 2.关闭

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
}

// 选项模式

type HTTPOption func(h *HTTPServer)

type HTTPServer struct {
	srv  *http.Server
	stop func() error
}

func WithHTTPServerStop(fn func() error) HTTPOption {
	return func(h *HTTPServer) {
		if fn == nil {
			fn = func() error {
				fmt.Println("111111111")
				quit := make(chan os.Signal)
				signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
				<-quit
				log.Println("shutdown Server....")

				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				if err := h.srv.Shutdown(ctx); err != nil {
					log.Fatal("server Shutdown", err)
				}
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
	h := &HTTPServer{}
	for _, opt := range opts {
		opt(h)
	}
	return h
}

// 接收请求转发请求
func (h *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}

func (h *HTTPServer) Start(addr string) error {
	h.srv = &http.Server{
		Addr:    addr,
		Handler: h,
	}
	return h.srv.ListenAndServe()
}

// 这里服务的关闭需要优雅的关闭
func (h *HTTPServer) Stop() error {
	return h.stop()
}
