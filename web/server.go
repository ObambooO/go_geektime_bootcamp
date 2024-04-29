package web

import (
	"net"
	"net/http"
)

type HandleFunc func(ctx Context)

// 确保HttpServer实现了Server接口
var _ Server = &HttpServer{}

type Server interface {
	// 实现2
	// Start()

	http.Handler
	Start(address string) error

	// AddRoute 路由注册功能
	/**
	 * method 是http方法
	 * path是路由
	 * handleFunc是业务逻辑
	 */
	//AddRoute(method string, path string, handleFunc HandleFunc)

	// AddRoute1 注册多个路由
	//AddRoute1(method string, path string, handleFunc ...HandleFunc)
}

//type HttpsServer struct {
//	HttpServer
//}

type HttpServer struct {
	// addr string // 创建的时候传递，而不是在Start的时候进行传递

	*Router
}

func NewHttpServer() *HttpServer {
	return &HttpServer{
		Router: newRouter(),
	}
}

//func (h *HttpServer) AddRoute1(method string, path string, handleFunc ...HandleFunc) {
//	//TODO implement me
//	panic("implement me")
//}

// AddRoute 路由注册功能
/**
 * method http方法
 * path 路由
 * handleFunc 业务逻辑
 */
//func (h *HttpServer) AddRoute(method string, path string, handleFunc HandleFunc) {
//	//TODO implement me
//	//panic("implement me")
//}

// Get get路由方法
func (h *HttpServer) Get(path string, handleFunc HandleFunc) {
	h.AddRoute(http.MethodGet, path, handleFunc)
}

func (h *HttpServer) Post(path string, handleFunc HandleFunc) {
	h.AddRoute(http.MethodPost, path, handleFunc)
}

func (h *HttpServer) Put(path string, handleFunc HandleFunc) {
	h.AddRoute(http.MethodPut, path, handleFunc)
}

func (h *HttpServer) Delete(path string, handleFunc HandleFunc) {
	h.AddRoute(http.MethodDelete, path, handleFunc)
}

// ServeHTTP 处理请求的入口
func (h *HttpServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	// 框架代码在这里
	ctx := &Context{
		Req:  request,
		Resp: writer,
	}
	// 查找路由，并且执行命中的业务逻辑
	h.serve(ctx)
}

func (h *HttpServer) serve(ctx *Context) {
	// 接下来是查找路由，并且执行命中的业务逻辑
}

func (h *HttpServer) Start(address string) error {
	l, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	// 在这里，可以让用户注册所谓的after start回调
	// 比如在这里往admin注册自己的这个实例
	// 在这里执行一些业务所需的前置条件

	return http.Serve(l, h)
}
