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
					chiledren: map[string]*node{
						"user": {
							part: "user",
							chiledren: map[string]*node{
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
