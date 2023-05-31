package trie

import "strings"

// Router维护一个根节点
type Router struct {
	root map[string]*node
}

// AddRouter 最开始我们手中有的数据肯定是类似这样的字符串
// /user/login hello
// /user/register world
// 就是将/user/login等字符串进行切割，分块保存到前缀树上

func (r *Router) AddRouter(pattern string, danode string) {
	if r.root == nil {
		r.root = make(map[string]*node)
	}
	root, ok := r.root["/"]
	// 创建根路由
	if !ok {
		root = &node{
			part: "/",
		}
		r.root["/"] = root
	}

	parts := strings.Split(strings.Trim(pattern, "/"), "/")
	for _, part := range parts {
		if part == "" {
			panic("pattern不符合格式...")
		}
		root.addNode(part)
	}
}

func (r *Router) GetRouter(pattern string) *node {
	return nil
}

type node struct {
	// part 当前节点的唯一标识
	part string
	// children 维护子节点的结构
	chiledren map[string]*node
	// data 当前节点需要保存的数据
	data string
}

// 这个节点有什么功能?
// 1.注册节点:新建一个node节点
// 2.查找节点

func (n *node) addNode(part string) *node {
	return nil
}

func (n *node) getNode(part string) *node {
	return nil
}
