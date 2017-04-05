package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSsidsHandler(t *testing.T) {

	ResourcesPath = "../static"
	SsidsFile = "../static/tests/ssids"

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

func TestReadSsidsFile(t *testing.T) {

	SsidsFile = "../static/tests/ssids"

	ssids, err := readSsidsFile()
	if err != nil {
		t.Errorf("Unexpected error reading ssids file: %v", err)
	}

	if len(ssids) != 4 {
		t.Error("Expected 4 elements in csv record")
	}

	set := make(map[string]bool)
	for _, v := range ssids {
		set[v] = true
	}

	if !set["mynetwork"] {
		t.Error("mynetwork value not found")
	}
	if !set["yournetwork"] {
		t.Error("yournetwork value not found")
	}
	if !set["hernetwork"] {
		t.Error("hernetwork value not found")
	}
	if !set["hisnetwork"] {
		t.Error("hisnetwork value not found")
	}
}

func TestReadSsidsFileWithOnlyOne(t *testing.T) {

	SsidsFile = "../static/tests/ssids_onlyonessid"

	ssids, err := readSsidsFile()
	if err != nil {
		t.Errorf("Unexpected error reading ssids file: %v", err)
	}

	if len(ssids) != 1 {
		t.Error("Expected 1 elements in csv record")
	}

	set := make(map[string]bool)
	for _, v := range ssids {
		set[v] = true
	}

	if !set["mynetwork"] {
		t.Error("mynetwork value not found")
	}
}

func TestReadEmptySsidsFile(t *testing.T) {

	SsidsFile = "../static/tests/ssids_empty"

	ssids, err := readSsidsFile()
	if err != nil {
		t.Errorf("Unexpected error reading ssids file: %v", err)
	}

	if len(ssids) != 0 {
		t.Error("Expected 0 elements in csv record")
	}
}
