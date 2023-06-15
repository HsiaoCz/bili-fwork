package bfwork

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRouterAdd(t *testing.T) {
	testCases := []struct {
		name    string
		method  string
		pattern string
		wantErr string
	}{
		{
			name:    "test1",
			method:  "GET",
			pattern: "/study/golang",
		},
		{
			name:    "test2",
			method:  "GET",
			pattern: "study/python",
			wantErr: "web:路由必须以 / 开头",
		},
		{
			name:    "test3",
			method:  "POST",
			pattern: "/study/java/",
			wantErr: "web:路由不能以 / 结尾",
		},
		{
			name:    "test4",
			method:  "GET",
			pattern: "study/rust/",
			wantErr: "web:路由必须以 / 开头",
		},
		{
			name:    "test5",
			method:  "GET",
			pattern: "/study/hello////",
			wantErr: "web:路由不能出现连续的 / ",
		},
	}

	r := newRouter()

	var mockHandleFunc HandleFunc = func(c *Context) {}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r.addRouter(tc.method, tc.pattern, mockHandleFunc)
			assert.PanicsWithError(t, tc.wantErr, func() {})
		})
	}
}

func TestRouterParmAdd(t *testing.T) {
	testCases := []struct {
		name     string
		method   string
		pattern  string
		wantBool bool
	}{
		{
			name:     "test1",
			method:   "GET",
			pattern:  "/study/:source",
			wantBool: false,
		},
		{
			name:     "test2",
			method:   "GET",
			pattern:  "/study/login1",
			wantBool: false,
		},
	}
	r := newRouter()
	var mockHandleFunc HandleFunc = func(c *Context) {}
	r.addRouter("GET", "/study/:course", mockHandleFunc)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			n, _, ok := r.getRouter(tc.method, tc.pattern)
			assert.Equal(t, tc.wantBool, ok)
			// 这里的n其实是一个参数路由
			// 参数路由有一个特点，就是它的part是以:开头
			assert.True(t, tc.wantBool, strings.HasPrefix(n.part, ":"))
			r.addRouter(tc.method, tc.pattern, mockHandleFunc)
		})
	}
}
