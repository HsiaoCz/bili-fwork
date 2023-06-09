package bfwork

import "strings"

// router路由树
/**
{
	"GET":node{},
	"POST":node{},
	"DETELE":node{},
	.....
}

**/
type router struct {
	trees map[string]*node
}

// addRouter 注册路由
// 注册路由还是有很多东西是需要考虑的
// 1、什么样的pattern是合法的?
// method的问题不需要考虑
// 1. pattern =""
// 2. pattern =/asdf/asd////
// 3. pattern = /
// 这里还有一个问题，pattern需不需要以 / 开头?
// pattern需不需要以 / 结尾?
// 这里定死，必须以 / 开头，而且不能以 /结尾
// 这里当pattern为空的时候
// 可以返回error吗？可以返回，但是不好
// 这里注册要么成功，要么失败，失败直接panic

func (r *router) addRouter(method string, pattern string, handleFunc HandleFunc) {
	if pattern == "" {
		panic("web:路由不能为空")
	}
	// TODO 如果是以根路由怎么办?
	if !strings.HasPrefix(pattern, "/") {
		panic("web:路由必须以 / 开头")
	}
	if strings.HasSuffix(pattern, "/") {
		panic("web:路由不能以 / 结尾")
	}

	// 获取根节点
	root, ok := r.trees[method]
	if !ok {
		// 根节点不存在
		// 1.创建根节点
		// 2.把根节点放到trees里面
		root = &node{
			part: "/",
		}
		r.trees[method] = root
	}

	// 切割pattern
	// /user/login => ["","user","login"]
	parts := strings.Split(pattern[1:], "/")
	for _, part := range parts {
		if part == "" {
			panic("web:路由不能出现连续的 / ")
		}
		root = root.addNode(part)
	}
	root.handleFunc = handleFunc
}

// getRouter 匹配路由
// method 需要校验吗？ method=gghhhhhh
// pattern 需要校验吗？
// pattern 一些简单的可以校验：就是说 /awbuildjs.ssdddd.asdd
// pattern = /user/login/  这种路由注册时是非法的，但是匹配时是合法的
// pattern = /user//login  非法的路由
func (r *router) getRouter(method string, pattern string) (*node, bool) {
	if pattern == "" {
		return nil, false
	}
	// TODO  / 这种路由怎么办？
	// 获取根节点
	root, ok := r.trees[method]
	if !ok {
		return nil, false
	}
	parts := strings.Split(strings.Trim(pattern, "/"), "/")
	for _, part := range parts {
		if part == "" {
			return nil, false
		}
		root = root.getNode(part)
		if root == nil {
			return nil, false
		}
	}
	return root, true
}

type node struct {
	part     string
	children map[string]*node
	// handleFunc 这里存的是当前节点上的视图函数
	handleFunc HandleFunc
}

func (n *node) addNode(part string) *node {
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
	// n的children属性都不存在
	if n.children == nil {
		return nil
	}
	// 正常思路
	child, ok := n.children[part]
	if !ok {
		return nil
	}
	return child
}
