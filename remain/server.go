package remain

import "net/http"

type server interface {
	// 这要求我们的engine实现多路复用功能
	http.Handler
	// Start() 启动服务
	Start(addr string) error
	// Stop() 停止服务
	Stop() error
}

type HTTPServer struct{}

func (h *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {}
func (h *HTTPServer) Start(addr string) error                          { return nil }
func (h *HTTPServer) Stop() error                                      { return nil }

func NewHTTPServer() *HTTPServer {
	return &HTTPServer{}
}
