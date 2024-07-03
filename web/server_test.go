//go:build e2e

package web

import (
	"fmt"
	"net/http"
	"testing"
)

func TestHttpServer_ServeHTTP(t *testing.T) {
	server := NewHttpServer()
	server.middlewares = []Middleware{
		func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				fmt.Println("第一个before")
				next(ctx)
				fmt.Println("第一个after")
			}
		},
		func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				fmt.Println("第二个before")
				next(ctx)
				fmt.Println("第二个after")
			}
		},
		func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				fmt.Println("第三个中断")
				//next(ctx)
				//fmt.Println("第二个after")
			}
		},
		func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				fmt.Println("第四个看不到")
				//next(ctx)
				//fmt.Println("第二个after")
			}
		},
	}
	server.ServeHTTP(nil, &http.Request{})
}

func TestServer(t *testing.T) {
	//var h Server = &HttpServer{}
	h := NewHttpServer()
	h.addRoute(http.MethodGet, "/user", func(ctx *Context) {
		fmt.Println("我是第一个方法")
		fmt.Println("我是第二个方法")
	})

	h.Get("/user2", func(ctx *Context) {
		fmt.Println("这是get方法")
	})

	h.Get("/order/detail", func(ctx *Context) {
		ctx.Resp.Write([]byte("hello, order detail"))
	})

	h.Get("/order/abc", func(ctx *Context) {
		ctx.Resp.Write([]byte(fmt.Sprintf("hello, %s", ctx.Req.URL.Path)))
	})

	h.Get("/order/*", func(ctx *Context) {
		ctx.Resp.Write([]byte("hello, order *"))
	})

	h.Post("/values/:id", func(ctx *Context) {
		id, err := ctx.PathValue("id").AsInt64()
		if err != nil {
			ctx.Resp.WriteHeader(422)
			ctx.Resp.Write([]byte("id 输入不对"))
			return
		}
		ctx.Resp.Write([]byte(fmt.Sprintf("id: %d", id)))
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

// 不需要提供，让他们自己装饰
// 线程安全的
//type SafeContext struct {
//	Context
//	mutex sync.RWMutex
//}
//
//func (c *SafeContext) RespJSONOK() error {
//	c.mutex.Lock()
//	defer c.mutex.Unlock()
//	return c.Context.RespJSONOK(http.StatusOK, val)
//}
