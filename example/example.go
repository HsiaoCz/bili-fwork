package main

import (
	bfwork "b-fwork"
	"log"
	"net/http"
)

func main() {
	h := bfwork.NewHTTP(bfwork.WithHTTPServerStop(nil))
	go func() {
		err := h.Start(":9090")
		// 这里不等于http.ErrServerClosed代表没有合法的关闭服务
		if err != nil && err != http.ErrServerClosed {
			log.Println("启动失败")
			log.Fatal(err)
		}
	}()
	err := h.Stop()
	if err != nil {
		log.Println("关闭失败")
		log.Fatal(err)
	}
}
