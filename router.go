package bfwork

import (
	"fmt"
	"strings"
)

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

func newRouter() *router {
	return &router{trees: map[string]*node{}}
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
// method = GET
// pattern = /
// handleFunc = HandleFunc()
// 意思是什么呢?就是说 / 节点绑定一个视图函数

func (r *router) addRouter(method string, pattern string, handleFunc HandleFunc) {
	// 打印一下注册的路由
	fmt.Printf("add router %s - %s\n", method, pattern)
	if pattern == "" {
		panic("web:路由不能为空")
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
	// TODO 如果是以根路由怎么办?

	if pattern == "/" {
		root.handleFunc = handleFunc
		return
	}

	if !strings.HasPrefix(pattern, "/") {
		panic("web:路由必须以 / 开头")
	}
	if strings.HasSuffix(pattern, "/") {
		panic("web:路由不能以 / 结尾")
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
func (r *router) getRouter(method string, pattern string) (*node, map[string]string, bool) {
	params := make(map[string]string)
	if pattern == "" {
		return nil, params, false
	}
	// TODO  / 这种路由怎么办？
	// 获取根节点

	root, ok := r.trees[method]
	if !ok {
		return nil, params, false
	}
	if pattern == "/" {
		return root, params, true
	}
	parts := strings.Split(strings.Trim(pattern, "/"), "/")
	for _, part := range parts {
		if part == "" {
			return nil, params, false
		}
		root = root.getNode(part)
		if root == nil {
			return nil, params, false
		}
		if strings.HasPrefix(root.part, ":") {
			params[root.part[1:]] = part
		}
		// /stufy/:course/action
	}
	return root, params, root.handleFunc != nil
}

type node struct {
	part string
	// 这个children 就是静态路由
	children map[string]*node
	// handleFunc 这里存的是当前节点上的视图函数
	handleFunc HandleFunc
	// 参数路由
	// 这里产生了一个疑问，为什么这里时一个纯node节点
	// /study/:source
	// /study/:programing
	// /study/golang 这个路由匹配哪个？
	// 这个根本匹配不了
	// 一个位置只能有一个动态参数，也就是一个占位符
	// 问题2：静态路由和动态路由的优先级
	// 注册的路由 /study/golang
	// 注册的路由 /study/:course
	// 请求的地址 /study/golang
	// 请求的静态路由的优先级要高于动态路由
	paramChildern *node
}

// addNode 这个方法是在服务启动前调用
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
	// n的children属性都不存在
	if n.children == nil {
		return nil
	}
	// 正常思路：先到静态路由中找
	child, ok := n.children[part]
	if !ok {
		// 到了这里说明没找到
		// 没找到，说明没有匹配到静态路由
		// 如果动态路由上有值 则返回动态路由
		if n.paramChildern != nil {
			return n.paramChildern
		}
		return nil
	}
	return child
}

// 路由分为动态路由和静态路由
// 静态路由：
// /study/golang
// /user/login
// /register

// 动态路由
// 1.参数路由
// /study/:course
// 这是注册的路由，匹配的时候可以匹配到/study/golang /study/python
// 但是类似这种 /study/golang/action 这种路由就匹配不到
// 2.通配符路由 贪婪匹配的
// /static/*filepath  注册时注册这种路由
// 匹配的时候，匹配到/static/css/style.css
// /static/js/index.js
// 3.正则路由
