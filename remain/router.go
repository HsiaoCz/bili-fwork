package remain

import "strings"

// 路由树
/**

{
	"GET":node{},
	"POST":node{},
	"DELETE":node{},
	.....
}

**/
// 每个请求方法维护一个路由树
type router struct {
	tress map[string]*node
}

func newRouter() *router {
	return &router{tress: map[string]*node{}}
}

// 树中的一个节点
type node struct {
	part string
	// children 静态路由
	children map[string]*node
	// handleFunc 存当前节点上的视图函数
	handleFunc HandleFunc
	// 参数路由
	// 一个位置只能有一个参数路由
	// 并且静态路由的优先级大于动态路由
	paramChildern *node
}

func (n *node) addNode(part string) *node {

	if strings.HasPrefix(part, ":") && n.paramChildern == nil {
		n.paramChildern = &node{part: part}
		return n.paramChildern
	}

	// 判断当前节点有没有children属性，就是是不是nil
	if n.children == nil {
		n.children = make(map[string]*node)
	}
	child, ok := n.children[part]
	if !ok {
		child = &node{
			part: part,
		}
		n.children[part] = child
	}
	return child
}

func (n *node) getNode(part string) *node {
	if n.children == nil {
		return nil
	}

	// 正常思路，先到静态路由中找
	child, ok := n.children[part]
	if !ok {
		if n.paramChildern != nil {
			return n.paramChildern
		}
		return nil
	}
	return child
}
