package trie

import (
	"errors"
	"strings"
)

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
		root = root.addNode(part)
	}
	// 循环结束后，此时的root是最底层的
	// 这时候，咱们得统一设置data的值
	root.data = danode
}

func (r *Router) GetRouter(pattern string) (*node, error) {
	root, ok := r.root["/"]
	// 创建根路由
	if !ok {
		return nil, errors.New("根节点不存在")
	}
	// 切割pattern
	// ["user","login"]
	parts := strings.Split(strings.Trim(pattern, "/"), "/")
	for _, part := range parts {
		if part == "" {
			return nil, errors.New("pattern格式不对")
		}
		root = root.getNode(part)
		if root == nil {
			return nil, errors.New("pattern 不存在")
		}
	}
	return root, nil
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
	// 判断当前节点有没有children属性，就是是不是nil
	if n.chiledren == nil {
		n.chiledren = make(map[string]*node)
	}
	child, ok := n.chiledren[part]
	if !ok {
		child = &node{
			part: part,
		}
		n.chiledren[part] = child
	}
	return child
}

func (n *node) getNode(part string) *node {
	// n的children属性都不存在
	if n.chiledren == nil {
		return nil
	}
	// 正常思路
	child, ok := n.chiledren[part]
	if !ok {
		return nil
	}
	return child
}
