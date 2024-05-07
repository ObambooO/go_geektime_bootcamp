package web

import "strings"

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

func (r *Router) AddRoute(method, path string, handleFunc HandleFunc) {
	// 首先找到树
	root, ok := r.trees[method]

	if !ok {
		// 说明没有根节点
		root = &node{
			path: "/",
		}
		r.trees[method] = root
	}

	// 去除最前面的/
	path = path[1:]
	// 切割path
	segments := strings.Split(path, "/")
	for _, segment := range segments {
		// 递归下去找准位置
		// 如果中途有节点不存在，则创建节点
		children := root.childOrCreate(segment)
		root = children
	}
	root.handleFunc = handleFunc
}

type node struct {
	// children 子节点
	children map[string]*node
	// path 节点路径
	path string
	// handleFunc 处理函数
	handleFunc HandleFunc
}

// childOf 返回segment对应的子节点，第一个值返回正确的子节点，第二个
func (n *node) childOrCreate(segment string) *node {
	if n.children == nil {
		n.children = map[string]*node{}
	}
	res, ok := n.children[segment]
	if !ok {
		// 新建一个
		res = &node{
			path: segment,
		}
		n.children[segment] = res
	}
	return res
}
