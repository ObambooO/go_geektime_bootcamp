package recover

import (
	"fmt"
	"testing"
	"web"
)

func TestMiddlewareBuilder_Build(t *testing.T) {
	builder := MiddlewareBuilder{
		StatusCode: 500,
		Data:       []byte("server error"),
		Log: func(ctx *web.Context) {
			fmt.Printf("panic 路径：%s", ctx.Req.URL.String())
		},
	}
	server := web.NewHttpServer(web.ServerWithMiddleware(builder.Build()))

	server.Get("/user", func(ctx *web.Context) {
		panic("user error")
	})

	server.Start(":8081")
}
