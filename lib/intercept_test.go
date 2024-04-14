package lib_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	greload "github.com/yamavol/greload/lib"
	"github.com/yamavol/greload/test/harness"
)

// ==================================================
// Middleware Test
// ==================================================
func Test_Middleware(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com", nil)
	w := httptest.NewRecorder()

	h := http.HandlerFunc(htmlHandler)
	r := greload.WithReload(h)
	r.ServeHTTP(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	var hasScript = strings.Contains(string(body), "<script>")
	var hasContentLength = len(resp.Header["Content-Length"]) > 0

	harness.IsEqual(t, resp.StatusCode, 200, "")
	harness.IsEqual(t, resp.Header["Content-Type"][0], "text/html; charset=utf-8", "check content type")
	harness.IsTrue(t, hasScript, "body should contain <script> tag")
	harness.IsFalse(t, hasContentLength, "has content length")
}

// ==================================================
// Other tests
// ==================================================
func Test_DetectContentType(t *testing.T) {
	harness.IsEqual(t,
		http.DetectContentType([]byte("<html></html>")),
		"text/html; charset=utf-8",
		"http.DetectContentType should identify HTML",
	)
}

// ==================================================
// dummy handlers
// ==================================================
func htmlHandler(w http.ResponseWriter, r *http.Request) {
	greload.WriteHtml(w, r, []byte("<html><body><h1>Hello!</h1></body></html>"))
}
