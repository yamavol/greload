package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/yamavol/greload"
)

func main() {

	go greload.WatchStart([]string{"."}, nil)

	r := http.FileServer(http.Dir("./"))
	r = greload.WithReload(r)

	serverPort := 4000
	server := &http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%v", serverPort),
		Handler: r,
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	go func() {
		<-interrupt
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		server.Shutdown(ctx)
	}()

	fmt.Println("http server is running on port", serverPort)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("Error:", err)
	}
}
