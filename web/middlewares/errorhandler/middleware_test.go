package errorhandler

import (
	"net/http"
	"testing"
	"web"
)

func TestMiddlewareBuilder_Build(t *testing.T) {
	builder := NewMiddlewareBuilder()
	builder.AddCode(http.StatusNotFound, []byte(`
<html>
  <body>
	<h1>404 Not Found</h1>
  </body>
</html>
`)).
		AddCode(http.StatusBadRequest, []byte(`
<html>
  <body>
	<h1>500 Inter Error</h1>
  </body>
</html>
`))
	server := web.NewHttpServer(web.ServerWithMiddleware(builder.Build()))
	server.Start(":8081")
}
