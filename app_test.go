package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func init() {
	logInit()

}

func TestHealthAPI(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HealthCheckHandler)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("HealthCheckHandler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	wantct := "application/json"
	if ct := rr.Header().Get("content-type"); ct != wantct {
		t.Errorf("HealthCheckHandler returned wrong content-type: got %v want %v",
			ct, wantct)
	}
}
