package accesslog

import (
	"fmt"
	"net/http"
	"testing"
	"web"
)

func TestMiddlewareBuilder(t *testing.T) {
	builder := MiddlewareBuilder{}
	middleware := builder.LogFunc(func(log string) {
		fmt.Println(log)
	}).Build()
	server := web.NewHttpServer(web.ServerWithMiddleware(middleware))
	server.Post("/api/*", func(ctx *web.Context) {
		fmt.Println("我是第一个方法")
	})
	req, err := http.NewRequest(http.MethodPost, "/api/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	server.ServeHTTP(nil, req)
}
