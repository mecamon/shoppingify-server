package auth

import (
	"net/http"
	"testing"
)

func TestRoutes(t *testing.T) {
	var i interface{}
	i = Routes()

	if _, ok := i.(http.Handler); !ok {
		t.Error("incorrect handler type")
	}
}
