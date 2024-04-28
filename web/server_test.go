package web

import (
	"fmt"
	"net/http"
	"testing"
)

func TestServer(t *testing.T) {
	//var h Server = &HttpServer{}
	h := &HttpServer{}
	h.AddRoute(http.MethodGet, "/user", func(ctx Context) {
		fmt.Println("我是第一个方法")
		fmt.Println("我是第二个方法")
	})

	h.Get("/user2", func(ctx Context) {
		fmt.Println("这是get方法")
	})

	// 注册多个不需要去管，让用户自己去处理
	//h.AddRoute1(http.MethodGet, "/user1", func(ctx Context) {
	//	fmt.Println("我是第一个方法")
	//}, func(ctx Context) {
	//	fmt.Println("我是第二个方法")
	//})

	h.Start(":8081")

	//go func() {
	//	http.ListenAndServe(":8080", h)
	//}()
}