package main

//run using go test ./...

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestServer(t *testing.T) {
	// Start a local HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, request *http.Request) {
		// Test request parameters
		equalto(t, request.URL.String(), webpgName)
		// Send response to be tested
		rw.Write([]byte(`OK`))
	}))
	defer server.Close()

	// Use Client & URL from our local test server
	body, err := http.Get(webpgName)

	notnil(t, err)
	equalto(t, []byte("OK"), body)

}

// fails the test if an err is not nil.
func notnil(tb testing.TB, err error) {
	if err != nil {
		tb.FailNow()
	}
}

// fails the test if exp is not equal to the interface.
func equalto(tb testing.TB, exp, i interface{}) {
	if !reflect.DeepEqual(exp, i) {
		tb.FailNow()
	}
}
