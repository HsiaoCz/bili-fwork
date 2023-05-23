package main

import (
	bfwork "b-fwork"
	"fmt"
	"log"
	"net/http"
)

func main() {
	h := bfwork.NewHTTP(bfwork.WithHTTPServerStop(nil))
	h.GET("/user", Login)
	// go func() {
	// 	err := h.Start(":9090")
	// 	// 这里不等于http.ErrServerClosed代表没有合法的关闭服务
	// 	if err != nil && err != http.ErrServerClosed {
	// 		log.Println("启动失败")
	// 		log.Fatal(err)
	// 	}
	// }()
	// err := h.Stop()
	// if err != nil {
	// 	log.Println("关闭失败")
	// 	log.Fatal(err)
	// }
	err := h.Start(":9091")
	if err != nil {
		log.Fatal(err)
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello")
}
