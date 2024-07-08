//go:build e2e

package accesslog

import (
	"fmt"
	"testing"
	"web"
)

func TestMiddlewareBuilderE2E(t *testing.T) {
	builder := MiddlewareBuilder{}
	middleware := builder.LogFunc(func(log string) {
		fmt.Println(log)
	}).Build()
	server := web.NewHttpServer(web.ServerWithMiddleware(middleware))
	server.Post("/api/*", func(ctx *web.Context) {
		fmt.Println("我是第一个方法")
	})
	server.Start(":8081")
}
