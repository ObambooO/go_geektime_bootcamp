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
			path:   "/",
		},
		{
			method: http.MethodGet,
			path:   "/user",
		},
		{
			method: http.MethodGet,
			path:   "/user/home",
		},
		{
			method: http.MethodGet,
			path:   "/order/detail",
		},
		{
			method: http.MethodGet,
			path:   "/order/*",
		},
		// 与正则冲突
		//{
		//	method: http.MethodGet,
		//	path:   "/order/detail/:id",
		//},
		{
			method: http.MethodPost,
			path:   "/user",
		},
		{
			method: http.MethodPost,
			path:   "/user/account",
		},
		{
			method: http.MethodGet,
			path:   "/*",
		},
		{
			method: http.MethodGet,
			path:   "/*/*",
		},
		{
			method: http.MethodGet,
			path:   "/*/abc",
		},
		{
			method: http.MethodGet,
			path:   "/*/abc/*",
		},
		// 目前暴力切无法处理下面的情况，需要校验
		//{
		//	method: http.MethodGet,
		//	path:   "login",
		////	path: "login/////",
		//},
		{
			method: http.MethodGet,
			path:   "/order/detail/:id(^[0-9]+$)",
		},
		{
			method: http.MethodGet,
			path:   "/order/detail/:id(^[0-9]+$)/:name(^[a-zA-Z]+$)",
		},
	}

	var mockHandler HandleFunc = func(ctx *Context) {}
	r := newRouter()
	for _, route := range testRoutes {
		r.addRoute(route.method, route.path, mockHandler)
	}

	// 断言路由树和预期的一致
	wantRouter := &Router{
		trees: map[string]*node{
			http.MethodGet: &node{
				path:       "/",
				handleFunc: mockHandler,
				children: map[string]*node{
					"user": &node{
						path:       "user",
						handleFunc: mockHandler,
						children: map[string]*node{
							"home": &node{
								path:       "home",
								handleFunc: mockHandler,
							},
						},
					},
					"order": &node{
						path: "order",
						children: map[string]*node{
							"detail": &node{
								path:       "detail",
								handleFunc: mockHandler,
								paramChild: &node{
									path:       ":id",
									handleFunc: mockHandler,
									regexpPath: "^[0-9]+$",
									paramChild: &node{
										path:       ":name",
										handleFunc: mockHandler,
										regexpPath: "^[a-zA-Z]+$",
									},
								},
							},
						},
						startChild: &node{
							path:       "*",
							handleFunc: mockHandler,
						},
					},
				},
				startChild: &node{
					path: "*",
					children: map[string]*node{
						"abc": &node{
							path: "abc",
							startChild: &node{
								path:       "*",
								handleFunc: mockHandler,
							},
						},
					},
				},
			},
			http.MethodPost: &node{
				path: "/",
				children: map[string]*node{
					"user": &node{
						path:       "user",
						handleFunc: mockHandler,
						children: map[string]*node{
							"account": &node{
								path:       "account",
								handleFunc: mockHandler,
							},
						},
					},
				},
			},
		},
	}

	msg, ok := wantRouter.equal(&r)
	assert.True(t, ok, msg)

	r = newRouter()
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "", mockHandler)
	}, "路径不能为空字符串")
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "user/login", mockHandler)
	}, "路径必须以/开头")
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "/user/login/", mockHandler)
	}, "路径不能以/结尾")
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "/user//login", mockHandler)
	}, "不能有连续的//")

	r = newRouter()
	r.addRoute(http.MethodGet, "/", mockHandler)
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "/", mockHandler)
	}, "路由冲突，重复注册[/]")

	r = newRouter()
	r.addRoute(http.MethodGet, "/a/b/c", mockHandler)
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "/a/b/c", mockHandler)
	}, "路由冲突，重复注册[/a/b/c]")

	// 校验method方法，可将addRoute方法改成私有的避免
	// 校验mockHandler是否为nil，传nil相当于没注册，不需要校验
	//r.addRoute("aaa", "/a/b/c", mockHandler)

	r = newRouter()
	r.addRoute(http.MethodGet, "/a/*", mockHandler)
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "/a/:id", nil)
	}, "web: 路由冲突，同时存在通配符和路径参数，已有通配符匹配")

	r = newRouter()
	r.addRoute(http.MethodGet, "/a/:id", mockHandler)
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "/a/*", nil)
	}, "web: 路由冲突，同时存在通配符和路径参数，已有参数匹配")

	r = newRouter()
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "/a/:id((", nil)
	}, "web: 路由存在不完整的((括号信息")
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

	if n.startChild != nil {
		msg, ok := n.startChild.equal(y.startChild)
		if !ok {
			return msg, ok
		}
	}

	if n.paramChild != nil {
		msg, ok := n.paramChild.equal(y.paramChild)
		if !ok {
			return msg, ok
		}
		// 递归校验
		if n.paramChild.paramChild != nil {
			msg, ok := n.paramChild.paramChild.equal(y.paramChild.paramChild)
			if !ok {
				return msg, ok
			}
		}
	}

	if n.regexpPath != y.regexpPath {
		return "正则表达式不相等", false
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

func TestRouter_findRoute(t *testing.T) {
	testRoutes := []struct {
		method string
		path   string
	}{
		{
			method: http.MethodGet,
			path:   "/",
		},
		{
			method: http.MethodGet,
			path:   "/user",
		},
		{
			method: http.MethodGet,
			path:   "/user/home",
		},
		{
			method: http.MethodGet,
			path:   "/order/detail",
		},
		{
			method: http.MethodGet,
			path:   "/order/*",
		},
		{
			method: http.MethodPost,
			path:   "/user",
		},
		{
			method: http.MethodPost,
			path:   "/user/account",
		},
		{
			method: http.MethodDelete,
			path:   "/order/detail",
		},
		{
			method: http.MethodDelete,
			path:   "/",
		},
		{
			method: http.MethodPost,
			path:   "/login/:username",
		},
		{
			method: http.MethodDelete,
			path:   "/*",
		},
		//{
		//	method: http.MethodGet,
		//	path:   "/order/detail/a/c/v/e",
		//},
		{
			method: http.MethodGet,
			path:   "/order/detail/:id(^[0-9]$+)",
		},
		{
			method: http.MethodGet,
			path:   "/order/detail/:id(^[0-9]$+)/name",
		},
	}

	r := newRouter()
	var mockHandler HandleFunc = func(ctx *Context) {}
	for _, route := range testRoutes {
		r.addRoute(route.method, route.path, mockHandler)
	}

	testCases := []struct {
		name      string
		method    string
		path      string
		wantFound bool
		matchInfo *matchInfo
	}{
		{
			// 根节点
			name:      "root",
			method:    http.MethodDelete,
			path:      "/",
			wantFound: true,
			matchInfo: &matchInfo{
				n: &node{
					handleFunc: mockHandler,
					path:       "/",
					children: map[string]*node{
						"order": &node{
							path: "order",
							children: map[string]*node{
								"detail": &node{
									handleFunc: mockHandler,
									path:       "detail",
								},
							},
						},
					},
					startChild: &node{
						path:       "*",
						handleFunc: mockHandler,
					},
				},
			},
		},
		{
			// 方法都不存在
			name:      "method no found",
			method:    http.MethodConnect,
			path:      "/order/detail",
			wantFound: false,
			matchInfo: &matchInfo{
				n: &node{
					handleFunc: mockHandler,
					path:       "detail",
				},
			},
		},
		{
			name:      "order start",
			method:    http.MethodGet,
			path:      "/order/abc",
			wantFound: true,
			matchInfo: &matchInfo{
				n: &node{
					handleFunc: mockHandler,
					path:       "*",
				},
			},
		},
		{
			// 命中但没有handler
			name:   "order",
			method: http.MethodGet,
			path:   "/order",
			// 这里true随方法里面的root, true后面而变更
			wantFound: true,
			matchInfo: &matchInfo{
				n: &node{
					//handleFunc: mockHandler,
					path: "order",
					children: map[string]*node{
						"detail": &node{
							handleFunc: mockHandler,
							path:       "detail",
						},
					},
				},
			},
		},
		{
			// username路径参数匹配
			name:      "login username",
			method:    http.MethodPost,
			path:      "/login/熊二",
			wantFound: true,
			matchInfo: &matchInfo{
				n: &node{
					path:       ":username",
					handleFunc: mockHandler,
				},
				pathParams: map[string]string{
					"username": "熊二",
				},
			},
		},
		// 下面当同时存在通配符和路径参数时，如果/api/detail和/api/*，则/api/detail会被匹配，另一个不会匹配
		{
			// 通配符多行匹配
			name:      "通配符多行匹配",
			method:    http.MethodGet,
			path:      "/order/a/v/c/e",
			wantFound: true,
			matchInfo: &matchInfo{
				n: &node{
					handleFunc: mockHandler,
					path:       "*",
				},
			},
		},
		{
			// 完全命中
			name:      "order detail",
			method:    http.MethodGet,
			path:      "/order/detail",
			wantFound: true,
			matchInfo: &matchInfo{
				n: &node{
					handleFunc: mockHandler,
					path:       "detail",
				},
			},
		},
		{
			name:      "通配符仅匹配一个",
			method:    http.MethodGet,
			path:      "/order/a",
			wantFound: true,
			matchInfo: &matchInfo{
				n: &node{
					handleFunc: mockHandler,
					path:       "*",
				},
			},
		},
		{
			name:      "匹配正则",
			method:    http.MethodGet,
			path:      "/order/detail/2",
			wantFound: true,
			matchInfo: &matchInfo{
				pathParams: map[string]string{
					"id": "2",
				},
				n: &node{
					handleFunc: mockHandler,
					path:       ":id",
					regexpPath: "^[0-9]$+",
					children: map[string]*node{
						"name": &node{
							path:       "name",
							handleFunc: mockHandler,
						},
					},
				},
			},
		},
		{
			name:      "匹配正则失败",
			method:    http.MethodGet,
			path:      "/order/detail/ssss",
			wantFound: false,
		},
		{
			name:      "匹配更深一级",
			method:    http.MethodGet,
			path:      "/order/detail/1/name",
			wantFound: true,
			matchInfo: &matchInfo{
				pathParams: map[string]string{
					"id": "1",
				},
				n: &node{
					handleFunc: mockHandler,
					path:       "name",
					regexpPath: "",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			info, found := r.findRoute(tc.method, tc.path)
			assert.Equal(t, tc.wantFound, found)
			if !found {
				return
			}
			assert.Equal(t, tc.matchInfo.pathParams, info.pathParams)
			msg, ok := tc.matchInfo.n.equal(info.n)
			assert.True(t, ok, msg)
		})

	}
}
