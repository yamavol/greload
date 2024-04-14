package http_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yamavol/greload/test/harness"
)

// ==================================================
// Handler Test (sample)
// ==================================================
func Test_HelloHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com", nil)
	w := httptest.NewRecorder()
	helloHandler(w, req)
	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	harness.IsEqual(t, resp.StatusCode, 200, "")
	harness.IsEqual(t, string(body), "Hello!", "")
	harness.IsEqual(t, resp.Header["Content-Type"][0], "text/plain; charset=utf-8", "check content type")
	// harness.IsEqual(t, resp.Header["Content-Length"][0], "6", "check content length")
}

func Test_HtmlHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com", nil)
	w := httptest.NewRecorder()
	htmlHandler(w, req)
	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	harness.IsEqual(t, resp.StatusCode, 200, "")
	harness.IsEqual(t, string(body), "<html><body><h1>Hello!</h1></body></html>", "")
	harness.IsEqual(t, resp.Header["Content-Type"][0], "text/html; charset=utf-8", "check content type")
}

func Test_EchoHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com?a=1", nil)
	w := httptest.NewRecorder()
	echoHandler(w, req)
	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	harness.IsEqual(t, string(body), "a: 1\n", "query param should be echoed")
}

// ==================================================
// dummy handlers
// ==================================================
func helloHandler(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("Hello!"))
}

func htmlHandler(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("<html><body><h1>Hello!</h1></body></html>"))
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	for key, values := range queryParams {
		for _, value := range values {
			fmt.Fprintf(w, "%s: %s\n", key, value)
		}
	}
}
