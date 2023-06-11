package bfwork

import (
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
