package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSsidsHandler(t *testing.T) {

	ResourcesPath = "../static"

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/", nil)
	http.HandlerFunc(SsidsHandler).ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got: %d", http.StatusOK, w.Code)
	}

	if !strings.Contains(w.Header().Get("Content-Type"), "text/html") {
		t.Error("Response content type is not expected text/html")
	}
}

func TestConnectHandler(t *testing.T) {

	ResourcesPath = "../static"

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/connect", nil)
	http.HandlerFunc(SsidsHandler).ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got: %d", http.StatusOK, w.Code)
	}

	if !strings.Contains(w.Header().Get("Content-Type"), "text/html") {
		t.Error("Response content type is not expected text/html")
	}
}

func TestInvalidTemplateHandler(t *testing.T) {

	ResourcesPath = "/invalidpath"

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/", nil)
	http.HandlerFunc(SsidsHandler).ServeHTTP(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got: %d", http.StatusInternalServerError, w.Code)
	}
}
