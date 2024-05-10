package web

import (
	"fmt"
	"strings"
)

// Router 用来支持路由树的操作
type Router struct {
	// trees http method => 路由树根节点
	trees map[string]*node
}

func newRouter() Router {
	return Router{
		trees: make(map[string]*node),
	}
}

// addRoute path必须以/开头，不能以/结尾，中间也不能有连续的//
func (r *Router) addRoute(method, path string, handleFunc HandleFunc) {
	if path == "" {
		panic("路径不能为空字符串")
	}

	// 首先找到树
	root, ok := r.trees[method]

	if !ok {
		// 说明没有根节点
		root = &node{
			path: "/",
		}
		r.trees[method] = root
	}
	// 开头不能没有/
	if path[0] != '/' {
		panic("路径必须以/开头")
	}

	if path == "/" {
		if root.handleFunc != nil {
			panic("路由冲突，重复注册[/]")
		}
		root.handleFunc = handleFunc
		return
	}

	// 结尾
	if path[len(path)-1] == '/' {
		panic("路径不能以/结尾")
	}

	// 去除最前面的/
	path = path[1:]
	// 切割path
	segments := strings.Split(path, "/")
	for _, segment := range segments {
		if segment == "" {
			panic("不能有连续的//")
		}
		// 递归下去找准位置
		// 如果中途有节点不存在，则创建节点
		children := root.childOrCreate(segment)
		root = children
	}
	if root.handleFunc != nil {
		panic(fmt.Sprintf("路由冲突，重复注册[%s]", path))
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

func (n *node) childOf(path string) (*node, bool) {
	if n.children == nil {
		return nil, false
	}
	child, ok := n.children[path]
	return child, ok
}

func (r *Router) findRoute(method string, path string) (*node, bool) {
	// 基本上是沿着树深度遍历
	root, ok := r.trees[method]

	if !ok {
		return nil, false
	}

	if path == "/" {
		return root, true
	}

	// 把前置和后置的/都去掉
	path = strings.Trim(path, "/")
	segments := strings.Split(path, "/")
	for _, segment := range segments {
		child, found := root.childOf(segment)
		if !found {
			return nil, false
		}
		root = child
	}
	// 代表有节点，且节点有注册handler，写true则不一定有
	//return root, root.handleFunc != nil
	return root, true
}
