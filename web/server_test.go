package web

import (
	"net/http"
	"testing"
)

func TestServer(t *testing.T) {
	var h Server = &HttpServer{}
	h.Start(":8081")

	go func() {
		http.ListenAndServe(":8080", h)
	}()
}
