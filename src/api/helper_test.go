package api

import (
	"net/http/httptest"
	"testing"
)

func setupTest() (*Api, *httptest.ResponseRecorder) {
	a := New(false)
	w := httptest.NewRecorder()
	return a, w
}

func testCheckStatusCode(t *testing.T, w *httptest.ResponseRecorder, statusCode int) {
	if w.Code != statusCode {
		t.Errorf("Expected response status code %d, got %d", statusCode, w.Code)
	}
}

func testCheckBody(t *testing.T, w *httptest.ResponseRecorder, body string) {
	if w.Body.String() != body {
		t.Errorf("Expected response body:\n\n%s\n\ngot:\n\n%s", body, w.Body.String())
	}
}
