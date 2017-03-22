// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2017 Canonical Ltd
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License version 3 as
// published by the Free Software Foundation.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package wifiap

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

// Testing http response comes with values
type mockTransportReturnsValues struct{}

func (mock *mockTransportReturnsValues) Do(req *http.Request) (*http.Response, error) {
	rawBody := `{"result":{"test1":"abc", "test2": "def"},"status":"OK","status-code":200,"type":"sync"}`

	response := http.Response{
		StatusCode: 200,
		Status:     "200 OK",
		Body:       ioutil.NopCloser(strings.NewReader(rawBody)),
	}

	return &response, nil
}

func TestResponseWithValues(t *testing.T) {
	mock := &mockTransportReturnsValues{}
	restClient := newRestClient(mock)
	response, _ := restClient.sendHTTPRequest("uri", "GET", nil)

	if len(response.Result) != 2 {
		t.Errorf("response length is %v when expected 2", len(response.Result))
	}

	if response.Result["test1"] != "abc" {
		t.Error("Content for key 'test1' not valid")
	}

	if response.Result["test2"] != "def" {
		t.Error("Content for key 'test2' not valid")
	}
}

type mockTransportReturnsNoValue struct{}

func (mock *mockTransportReturnsNoValue) Do(req *http.Request) (*http.Response, error) {
	rawBody := `{"result":{},"status":"OK","status-code":200,"type":"sync"}`

	response := http.Response{
		StatusCode: 200,
		Status:     "200 OK",
		Body:       ioutil.NopCloser(strings.NewReader(rawBody)),
	}

	return &response, nil
}

func TestResponseWithoutValues(t *testing.T) {
	mock := &mockTransportReturnsNoValue{}
	restClient := newRestClient(mock)
	response, _ := restClient.sendHTTPRequest("uri", "GET", nil)

	if len(response.Result) > 0 {
		t.Errorf("response length is %v when expected 0", len(response.Result))
	}
}

// Testing http response is an error
type mockTransportReturnsError struct{}

func (mock *mockTransportReturnsError) Do(req *http.Request) (*http.Response, error) {
	return nil, fmt.Errorf("Failed: A random error message")
}

func TestErrorResponse(t *testing.T) {
	mock := &mockTransportReturnsError{}
	restClient := newRestClient(mock)
	_, err := restClient.sendHTTPRequest("uri", "GET", nil)

	if err == nil {
		t.Error("Expected an error, but got no response error")
	}

	if err.Error() != "Failed: A random error message" {
		t.Error("Got wrong error message")
	}
}
