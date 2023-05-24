# bili-fwork

这个框架还是有点类似于 gin

这个框架的核心还是在路由树

## 1、server 抽象

这里的 server 接口约束的是类似于 gin 框架里的 engine 的方法集

```go
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
```

这里 stop 考虑优雅关闭
这里通过选项模式，来给用户提供一个口子，让用户自定义优雅关闭的方法
用户如果不用，我们给出一个默认的优雅的关闭的方法

```go

type HTTPOption func(h *HTTPServer)

func NewHTTP(opts ...HTTPOption) *HTTPServer {
	h := &HTTPServer{
		routers: map[string]HandleFunc{},
	}
	for _, opt := range opts {
		opt(h)
	}
	return h
}

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

// 这里服务的关闭需要优雅的关闭
func (h *HTTPServer) Stop() error {
	return h.stop()
}
```

在使用的时候，start()需要另开一个线程

```go
	 go func() {
	 	err := h.Start(":9090")
	 	// 这里不等于http.ErrServerClosed代表没有合法的关闭服务
	 	if err != nil && err != http.ErrServerClosed {
	 		log.Println("启动失败")
	 		log.Fatal(err)
	 	}
	 }()
```

http.ErrServerClosed 代表正常关闭的错误

## 2、context

context 构造请求的上下文

```go

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
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		W:       w,
		R:       r,
		Method:  r.Method,
		Pattern: r.URL.Path,
	}
}
```

## 3、前缀树

任意一个节点和它的子节点拥有相同的前缀
