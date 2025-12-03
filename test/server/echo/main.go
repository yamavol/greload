package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"time"
)

func writeHeaders(w http.ResponseWriter, r *http.Request) {
	keys := make([]string, 0, len(r.Header))

	for k := range r.Header {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		for _, h := range r.Header[k] {
			fmt.Fprintf(w, "%v: %v\n", k, h)
		}
	}
}

func writeBody(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Body:\n%s", body)
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	w.Write([]byte(r.Proto + " "))
	w.Write([]byte(r.Method + " "))
	w.Write([]byte(r.RequestURI + "\n\n"))
	writeHeaders(w, r)
	w.Write([]byte("\n"))
	writeBody(w, r)
}

func main() {
	serverPort := 5000
	server := &http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%v", serverPort),
		Handler: http.HandlerFunc(echoHandler),
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	go func() {
		<-interrupt
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		server.Shutdown(ctx)
	}()

	fmt.Println("echo server is running on port", serverPort)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("Error:", err)
	}
}
