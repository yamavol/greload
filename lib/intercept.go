package lib

import (
	"bytes"
	"context"
	_ "embed"
	"io"
	"net/http"
	"strconv"
	"strings"
)

//go:embed reload-client.js
var reloadCode string
var injectHtml = "\n<script>\n" + reloadCode + "\n</script>\n"

const keyContentType = "Content-Type"
const keyContentLength = "Content-Length"

// InjectReload is a responseModifier function for ReverseProxy.ModifyResponse
// which injects reload script into HTML response.
func responseModifier(wsport int) func(*http.Response) error {

	var scriptHtml = []byte(strings.ReplaceAll(injectHtml, "9765", strconv.Itoa(wsport)))

	return func(resp *http.Response) error {
		if !strings.HasPrefix(resp.Header.Get(keyContentType), "text/html") {
			return nil
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		if err = resp.Body.Close(); err != nil {
			return err
		}

		body = append(body, scriptHtml...)
		resp.Body = io.NopCloser(bytes.NewReader(body))
		resp.ContentLength = int64(len(body))
		resp.Header.Set(keyContentLength, strconv.Itoa(len(body)))
		return nil
	}
}

// ============================================================
// greload (Manual Code Injection)
//
// The code below were initially part of greload (as library).
// It was abandoned because of the design failure. GoLang's net/http
// does not allow to rewrite the response body once written.
//
// Integrating these code is not a good idea because it introduces
// unnecessary dependencies into your program, but useful in some cases.
type greloadKeyType struct{}

var contextKey = greloadKeyType{}

// returns HTML with reload script injected, only if enabled
func InjectHtml(r *http.Request, html []byte) []byte {
	if r.Context().Value(contextKey) == true {
		return append(html, injectHtml...)
	} else {
		return html
	}
}

// Write Html writes HTML with reload script injected, only if enabled
func WriteHtml(w http.ResponseWriter, r *http.Request, html []byte) {
	w.Write(InjectHtml(r, html))
}

func WithReload(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := context.WithValue(r.Context(), contextKey, true)
		r = r.WithContext(c)
		next.ServeHTTP(w, r)
	}
}
