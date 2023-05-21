package bfwork

import "testing"

func TestServer(t *testing.T) {
	s := NewHTTP(WithHTTPServerStop(nil))

	go func() {
		if err := s.Start(":9090"); err != nil {
			t.Fail()
		}
	}()
	if err := s.Stop(); err != nil {
		t.Fail()
	}
}
