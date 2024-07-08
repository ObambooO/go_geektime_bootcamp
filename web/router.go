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
		panic("web: 路由是空字符串")
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
		panic("web: 路由必须以 / 开头")
	}

	if path == "/" {
		if root.handleFunc != nil {
			panic("web: 路由冲突[/]")
		}
		root.handleFunc = handleFunc
		root.route = "/"
		return
	}

	// 结尾
	if path[len(path)-1] == '/' {
		panic("web: 路由不能以 / 结尾")
	}

	// 去除最前面的/
	segments := strings.Split(path[1:], "/")

	//segments := listSegments(path)

	for _, segment := range segments {
		if segment == "" {
			panic(fmt.Sprintf("web: 非法路由。不允许使用 //a/b, /a//b 之类的路由, [%s]", path))
		}
		root = root.childOrCreate(segment)
	}
	if root.handleFunc != nil {
		panic(fmt.Sprintf("web: 路由冲突[%s]", path))
	}
	root.handleFunc = handleFunc
	root.route = path
}

// 目的，为了通配符的匹配
type nodeType int

const (
	// 静态路由
	nodeTypeStatic = iota
	// 正则路由
	nodeTypeReg
	// 路径参数路由
	nodeTypeParam
	// 通配符路由
	nodeTypeAny
)

type node struct {
	// children 子节点
	children map[string]*node
	// path 节点路径
	path  string
	route string

	// handleFunc 处理函数
	handleFunc HandleFunc

	// 通配符匹配
	starChild *node

	// 路径参数匹配
	paramChild *node
	// 正则路由和参数路由都会使用这个字段
	paramName string

	// 正则表达式
	regChild *node
	regExpr  *regexp.Regexp
	typ      nodeType
}

// childOrCreate 返回segment对应的子节点，第一个值返回正确的子节点，第二个
func (n *node) childOrCreate(segment string) *node {

	// 检验有没有重复注册
	if segment == "*" {
		if n.paramChild != nil {
			panic(fmt.Sprintf("web: 非法路由，已有路径参数路由。不允许同时注册通配符路由和参数路由 [%s]", segment))
		}
		if n.regChild != nil {
			panic(fmt.Sprintf("web: 非法路由，已有正则路由。不允许同时注册通配符路由和正则路由 [%s]", segment))
		}
		if n.starChild == nil {
			n.starChild = &node{
				path: segment,
				typ:  nodeTypeAny,
			}
		}
		return n.starChild
	}

	// 以：开头，需要进一步解析，判定是参数路由还是正则路由
	if segment[0] == ':' {
		paramName, expr, isReg := n.parseParam(segment)

		if isReg {
			return n.childOrCreateReg(segment, expr, paramName)
		}
		return n.childOrCreateParam(segment, paramName)
	}

	if n.children == nil {
		n.children = map[string]*node{}
	}
	res, ok := n.children[segment]
	if !ok {
		// 新建一个
		res = &node{
			path: segment,
			typ:  nodeTypeStatic,
		}
		n.children[segment] = res
	}
	return res
}

// childOf 优先考虑静态匹配，匹配不上，再考虑通配符匹配
// 第一个返回值是子节点，第二个是标记是否为路径参数
func (n *node) childOf(path string) (*node, bool) {
	if n.children == nil {
		n.childOfNonStatic(path)
	}
	child, ok := n.children[path]
	if !ok {
		return n.childOfNonStatic(path)
	}
	return child, ok
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

	mi := &matchInfo{}

	for _, segment := range segments {
		var child *node
		child, ok = root.childOf(segment)
		if !ok {
			if root.typ == nodeTypeAny {
				mi.n = root
				return mi, true
			}
			return nil, false
		}

		if child.paramName != "" {
			mi.addValue(child.paramName, segment)
		}
		root = child
	}

	mi.n = root
	return mi, true
}

type matchInfo struct {
	n          *node
	pathParams map[string]string
}

// parseParam 用于解析判定是不是正则表达式
// 第一个返回是参数名字
// 第二个返回值是正则表达式
// 第三个返回值为 true 则说明是正则路由
func (n *node) parseParam(path string) (string, string, bool) {
	// 去除:
	path = path[1:]
	segments := strings.SplitN(path, "(", 2)

	if len(segments) == 2 {
		expr := segments[1]
		if strings.HasSuffix(expr, ")") {
			return segments[0], expr[:len(expr)-1], true
		}
	}
	return path, "", false
}

// 创建子参数节点
func (n *node) childOrCreateParam(path string, paramName string) *node {
	if n.regChild != nil {
		panic(fmt.Sprintf("web: 非法路由，已有正则路由。不允许同时注册正则路由和参数路由 [%s]", path))
	}
	if n.starChild != nil {
		panic(fmt.Sprintf("web: 非法路由，已有通配符路由。不允许同时注册通配符路由和参数路由 [%s]", path))
	}
	if n.paramChild != nil {
		if n.paramChild.path != path {
			panic(fmt.Sprintf("web: 路由冲突，参数路由冲突，已有 %s，新注册 %s", n.paramChild.path, path))
		}
	} else {
		n.paramChild = &node{path: path, paramName: paramName, typ: nodeTypeParam}
	}
	return n.paramChild
}

// 创建正则节点
func (n *node) childOrCreateReg(path string, expr string, paramName string) *node {
	if n.starChild != nil {
		panic(fmt.Sprintf("web: 非法路由，已有通配符路由。不允许同时注册通配符路由和正则路由 [%s]", path))
	}
	if n.paramChild != nil {
		panic(fmt.Sprintf("web: 非法路由，已有路径参数路由，不允许同时注册正则路由和参数路由 [%s]", path))
	}
	if n.regChild != nil {
		if n.regChild.regExpr.String() != expr || n.paramName != paramName {
			panic(fmt.Sprintf("web: 路由冲突，正则路由冲突，已有 %s，新注册 %s", n.regChild.path, path))
		}
	} else {
		regExpr, err := regexp.Compile(expr)
		if err != nil {
			panic(fmt.Errorf("web: 正则表达式错误 %w", err))
		}
		n.regChild = &node{path: path, paramName: paramName, regExpr: regExpr, typ: nodeTypeReg}
	}
	return n.regChild
}

// childOfNonStatic 从非静态匹配的子节点里面查找
func (n *node) childOfNonStatic(path string) (*node, bool) {
	if n.regChild != nil {
		if n.regChild.regExpr.Match([]byte(path)) {
			return n.regChild, true
		}
	}
	if n.paramChild != nil {
		return n.paramChild, true
	}
	return n.starChild, n.starChild != nil
}

// 添加节点值
func (m *matchInfo) addValue(key string, value string) {
	if m.pathParams == nil {
		// 大多数情况下，参数路径只有一段
		m.pathParams = map[string]string{key: value}
	}
	m.pathParams[key] = value
}
