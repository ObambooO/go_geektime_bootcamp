package web

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
	"testing"
)

func TestRouter_AddRoute(t *testing.T) {
	// 构造路由树
	// 验证路由树
	testRoutes := []struct {
		method string
		path   string
	}{
		{
			method: http.MethodGet,
			path:   "/user/home",
		},
	}

	var mockHandler HandleFunc = func(ctx Context) {}
	r := newRouter()
	for _, route := range testRoutes {
		r.addRoute(route.method, route.path, mockHandler)
	}

	// 断言路由树和预期的一致
	wantRouter := &Router{
		trees: map[string]*node{
			http.MethodGet: &node{
				path: "/",
				children: map[string]*node{
					"user": &node{
						path: "user",
						children: map[string]*node{
							"home": &node{
								path:       "home",
								handleFunc: mockHandler,
							},
						},
					},
				},
			},
		},
	}

	msg, ok := wantRouter.equal(r)
	assert.True(t, ok, msg)

}

// 返回string是为了返回错误信息，帮助我们排查问题
// bool 代表是否真的相等
func (r *Router) equal(y *Router) (string, bool) {
	for k, v := range r.trees {
		dst, ok := y.trees[k]
		if !ok {
			return fmt.Sprintf("找不到对应的http method"), false
		}
		// v, dst 要相等
		msg, equal := v.equal(dst)
		if !equal {
			return msg, false

		}
	}
	return "", true
}

func (n *node) equal(y *node) (string, bool) {
	if n.path != y.path {
		return fmt.Sprintf("节点path不相等"), false
	}
	if len(n.children) != len(y.children) {
		return fmt.Sprintf("子节点数量不相等"), false
	}

	// 比较handler，reflect.ValueOf反射
	nHandler := reflect.ValueOf(n.handleFunc)
	yHandler := reflect.ValueOf(y.handleFunc)
	if nHandler != yHandler {
		return fmt.Sprintf("handler不相等"), false
	}

	for path, c := range n.children {
		dst, ok := y.children[path]
		if !ok {
			return fmt.Sprintf("找不到对应的子节点"), false
		}
		msg, ok := c.equal(dst)
		if !ok {
			return msg, false
		}
	}

	return "", true
}
