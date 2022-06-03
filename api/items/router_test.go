package items

import (
	"net/http"
	"testing"
)

func TestRoutes(t *testing.T) {
	var i interface{}
	i = Routes()

	if _, ok := i.(http.Handler); !ok {
		t.Error("error checking Routes() return type")
	}
}