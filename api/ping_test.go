package api

import (
	"net/http"
	"testing"
)

func TestPing(t *testing.T) {
	a, w := setupTest()

	req, _ := http.NewRequest("GET", "/ping", nil)
	a.engine.ServeHTTP(w, req)

	testCheckStatusCode(t, w, 200)
	testCheckBody(t, w, "pong")
}
