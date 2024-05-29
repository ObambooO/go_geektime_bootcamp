package web

import (
	"fmt"
	"regexp"
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

	stack := []rune{}
	bracketMap := map[rune]rune{
		')': '(',
	}

	// TODO 处理判定是否存在左右括号
	for _, character := range path {
		switch character {
		case '(':
			stack = append(stack, character)

		case ')':
			if len(stack) == 0 || stack[len(stack)-1] != bracketMap[character] {
				panic("web: 存在不完整的括号信息")
			}
			stack = stack[:len(stack)-1]
		}
	}

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

	// 通配符匹配
	startChild *node

	// 路径参数匹配
	paramChild *node

	// 正则参数
	regexpPath string
}

// childOrCreate 返回segment对应的子节点，第一个值返回正确的子节点，第二个
func (n *node) childOrCreate(segment string) *node {

	if segment[0] == ':' {
		// 检测是否有同时注册路径参数和通配符
		if n.startChild != nil {
			panic("web: 不允许同时注册路径参数和通配符匹配，已有通配符匹配")
		}
		n.paramChild = &node{
			path: segment,
		}
		return n.paramChild
	}

	// 检验有没有重复注册
	if segment == "*" {
		if n.paramChild != nil {
			panic("web: 不允许同时注册路径参数和通配符匹配，已有路径参数")
		}
		n.startChild = &node{
			path: segment,
		}
		return n.startChild
	}
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

// childOf 优先考虑静态匹配，匹配不上，再考虑通配符匹配
// 第一个返回值是子节点，第二个是标记是否为路径参数，第三个标记命中了没有
func (n *node) childOf(path string) (*node, bool, bool) {
	if n.children == nil {
		if n.paramChild != nil {
			return n.paramChild, true, true
		}
		return n.startChild, false, n.startChild != nil
	}
	child, ok := n.children[path]
	if !ok {
		if n.paramChild != nil {
			return n.paramChild, true, true
		}
		return n.startChild, false, n.startChild != nil
	}
	return child, false, ok
}

// findRoute 查找路由
func (r *Router) findRoute(method string, path string) (*matchInfo, bool) {
	// 基本上是沿着树深度遍历
	root, ok := r.trees[method]

	if !ok {
		return nil, false
	}

	if path == "/" {
		return &matchInfo{
			n: root,
		}, true
	}

	// 把前置和后置的/都去掉
	path = strings.Trim(path, "/")
	segments := strings.Split(path, "/")

	var pathParams map[string]string
	var starMatchInfo *matchInfo
	for _, segment := range segments {
		child, paramChild, found := root.childOf(segment)
		if !found {
			if starMatchInfo != nil {
				return starMatchInfo, true
			}
			return nil, false
		}

		// 命中了路径参数
		if paramChild {
			if pathParams == nil {
				pathParams = make(map[string]string)
			}
			// path是 :id这种形式，需要把:去掉
			pathParams[child.path[1:]] = segment
		}
		root = child
		if child.path == "*" {
			starMatchInfo = &matchInfo{
				n:          child,
				pathParams: pathParams,
			}

		}
	}
	// 代表有节点，且节点有注册handler，写true则不一定有
	//return root, root.handleFunc != nil
	return &matchInfo{
		n:          root,
		pathParams: pathParams,
	}, true
}

type matchInfo struct {
	n          *node
	pathParams map[string]string
}

func checkRegex(params string) bool {
	_, err := regexp.Compile(params)
	return err == nil
}
