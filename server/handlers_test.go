/*
 * Copyright (C) 2017 Canonical Ltd
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/CanonicalLtd/UCWifiConnect/utils"
)

func TestManagementHandler(t *testing.T) {

	ResourcesPath = "../static"
	SsidsFile := "../static/tests/ssids"
	utils.SetSsidsFile(SsidsFile)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/", nil)
	http.HandlerFunc(ManagementHandler).ServeHTTP(w, r)

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
	http.HandlerFunc(ManagementHandler).ServeHTTP(w, r)

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
	http.HandlerFunc(ManagementHandler).ServeHTTP(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got: %d", http.StatusInternalServerError, w.Code)
	}
}

func TestReadSsidsFile(t *testing.T) {

	SsidsFile := "../static/tests/ssids"
	utils.SetSsidsFile(SsidsFile)
	ssids, err := utils.ReadSsidsFile()
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

	SsidsFile := "../static/tests/ssids_onlyonessid"
	utils.SetSsidsFile(SsidsFile)
	ssids, err := utils.ReadSsidsFile()
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

	SsidsFile := "../static/tests/ssids_empty"
	utils.SetSsidsFile(SsidsFile)
	ssids, err := utils.ReadSsidsFile()
	if err != nil {
		t.Errorf("Unexpected error reading ssids file: %v", err)
	}

	if len(ssids) != 0 {
		t.Error("Expected 0 elements in csv record")
	}
}
