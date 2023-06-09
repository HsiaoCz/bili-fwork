package trie

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRouter_AddRouter(t *testing.T) {
	testCase := []struct {
		name       string
		parttern   string
		data       string
		wantRouter *Router
	}{
		{
			name:     "xxx",
			parttern: "/user/login",
			data:     "hello",
			wantRouter: &Router{map[string]*node{
				"/": {
					part: "/",
					children: map[string]*node{
						"user": {
							part: "user",
							children: map[string]*node{
								"login": {
									part: "login",
									data: "hello",
								},
							},
						},
					},
				},
			}},
		},
	}

	router := &Router{map[string]*node{
		"/": {
			part: "/",
		},
	}}
	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			router.AddRouter(tc.parttern, tc.data)
			assert.Equal(t, tc.wantRouter, router)
		})
	}
}

func TestRouter_GetRouter(t *testing.T) {
	testCase := []struct {
		// 测试的名字，随意给就好
		name string
		// 想要匹配的节点
		findparttern string
		// 想要返回的数据
		wantdata string
		// 理想中的错误
		wanterr error
	}{
		{
			name:         "success",
			findparttern: "/user/login",
			wantdata:     "hello",
		},
	}
	router := &Router{map[string]*node{
		"/": {
			part: "/",
		},
	}}
	router.AddRouter("/user/login", "hello")
	router.AddRouter("/user/register", "world")
	router.AddRouter("/study/golang", "Good")
	router.AddRouter("/study/python", "aaaa")
	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			n, err := router.GetRouter(tc.findparttern)
			assert.Equal(t, tc.wanterr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantdata, n.data)
		})
	}
}
