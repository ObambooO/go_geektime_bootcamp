package web

// Router 用来支持路由树的操作
type Router struct {
	// trees http method => 路由树根节点
	trees map[string]*node
}

func newRouter() *Router {
	return &Router{
		trees: make(map[string]*node),
	}
}

func (r *Router) addRoute(method, path string, handleFunc HandleFunc) {

}

type node struct {
	// children 子节点
	children map[string]*node
	// path 节点路径
	path string
	// handleFunc 处理函数
	handleFunc HandleFunc
}
