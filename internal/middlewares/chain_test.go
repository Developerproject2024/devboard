package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestChainAppliesMiddlewaresInOrder(t *testing.T) {
	order := []string{}
	middleware := func(name string) Middleware {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				order = append(order, name+"-before")
				next.ServeHTTP(writer, request)
				order = append(order, name+"-after")
			})
		}
	}

	handler := Chain(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		order = append(order, "handler")
	}), middleware("first"), middleware("second"))

	handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))

	expected := []string{"first-before", "second-before", "handler", "second-after", "first-after"}
	if len(order) != len(expected) {
		t.Fatalf("orden = %v; se esperaban %v", order, expected)
	}
	for i := range expected {
		if order[i] != expected[i] {
			t.Fatalf("orden = %v; se esperaban %v", order, expected)
		}
	}
}
