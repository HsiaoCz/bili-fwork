package remain

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

type server interface {
	// 这要求我们的engine实现多路复用功能
	http.Handler
	// Start() 启动服务
	Start(addr string) error
	// Stop() 停止服务
	Stop() error
}

// 这里定义一个视图函数
type HandleFunc func(w http.ResponseWriter, r *http.Request)

// 这里用选项模式添加一个优雅关机的功能

type HTTPOption func(h *HTTPServer)

func WithHTTPServerStop(fn func() error) HTTPOption {
	return func(h *HTTPServer) {
		if fn == nil {
			fn = func() error {
				fmt.Println("11111111")
				quit := make(chan os.Signal)
				// 如果匹配到中断信息，会将信号传到channel里面
				// 如果没有匹配到，channel会一直阻塞
				signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
				<-quit
				log.Println("shutdown Server ....")

				// 创建一个超时的上下文，在这里等待5秒钟
				// 等待任务执行完毕
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				if err := h.srv.Shutdown(ctx); err != nil {
					log.Fatal("server shutdown", err)
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

type HTTPServer struct {
	srv  *http.Server
	stop func() error
}

func (h *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {}
func (h *HTTPServer) Start(addr string) error {
	h.srv = &http.Server{
		Addr:    addr,
		Handler: h,
	}
	log.Println("the server is runing on port", addr)
	return h.srv.ListenAndServe()
}
func (h *HTTPServer) Stop() error { return h.stop() }

func NewHTTPServer(opts ...HTTPOption) *HTTPServer {
	h := &HTTPServer{}
	for _, opt := range opts {
		opt(h)
	}
	return h
}
