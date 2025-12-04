package lib

// ============================================================
// HTTP & Websocket server Methods
// ============================================================

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/mattn/go-ieproxy"
	"github.com/yamavol/greload/lib/internal"
	"github.com/yamavol/greload/log"

	"golang.org/x/net/websocket"
)

// ProxyServer is the libary's main server that handles both HTTP proxying
// and WebSocket connections.
type ProxyServer struct {
	options     ServerOptions
	connections map[*websocket.Conn]struct{}
	mu          sync.Mutex
	reloadReq   notifier
}

// Create a new instance of ProxyServer
func NewServer(options *ServerOptions) *ProxyServer {
	return &ProxyServer{
		options:     *options,
		connections: make(map[*websocket.Conn]struct{}),
		reloadReq:   *newNotifier(),
	}
}

// Start HTTP & WebSocket server, and listen to file change events
func (srv *ProxyServer) Start() {
	server := http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%v", srv.options.Port),
		Handler: http.HandlerFunc(serverHandler(srv)),
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	debounced, _ := internal.NewDebouncer(100 * time.Millisecond)

	go func() {
		<-interrupt
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		server.Shutdown(ctx)
	}()

	go func() {
		for range srv.reloadReq.ch {
			debounced(srv.handleReload)
		}
	}()

	log.Info("reload server is running on port", srv.options.Port)
	log.Info("redirecting access to", srv.options.Host.String())
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error("Error:", err)
	}
}

// Sends reload request to the connected websocket clients.
func (srv *ProxyServer) TriggerReload() {
	srv.reloadReq.Notify()
}

func serverHandler(srv *ProxyServer) func(w http.ResponseWriter, r *http.Request) {

	// reverse proxy server config
	rp := &httputil.ReverseProxy{}

	rp.ModifyResponse = internal.ResponseModifier(srv.options.Port)
	rp.Transport = &http.Transport{
		// use http_proxy (env) if set, otherwise use system proxy
		Proxy:             ieproxy.GetProxyFunc(),
		DisableKeepAlives: true,
	}
	rp.Rewrite = func(pr *httputil.ProxyRequest) {
		pr.SetXForwarded()
		pr.SetURL(srv.options.Host)
	}

	// websocket server config
	ws := websocket.Handler(func(conn *websocket.Conn) {
		srv.websockHandler(conn)
	})

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Upgrade") == "websocket" {
			ws.ServeHTTP(w, r)
		} else {
			rp.ServeHTTP(w, r)
		}
	}
}

func (ws *ProxyServer) websockHandler(conn *websocket.Conn) {
	defer conn.Close()

	ws.mu.Lock()
	ws.connections[conn] = struct{}{}
	ws.mu.Unlock()

	log.Debug("[ws] Client connected")

	var msg string
	for {
		err := websocket.Message.Receive(conn, &msg)
		if err != nil {
			ws.mu.Lock()
			delete(ws.connections, conn)
			ws.mu.Unlock()
			log.Debug("[ws] Client disconnected")
			break
		}
	}
}

func (srv *ProxyServer) handleReload() {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	for conn := range srv.connections {
		err := websocket.Message.Send(conn, "reload")
		if err != nil {
			delete(srv.connections, conn)
			log.Errorf("Error broadcasting message to client: %v", err)
		}
	}
}

// ============================================================
// notifier (Private)
// ============================================================

type notifier struct {
	ch chan struct{}
}

func newNotifier() *notifier {
	return &notifier{
		// create buffered channel of size 1
		// notify call is always unblocking.
		ch: make(chan struct{}, 1),
	}
}

// signals channel or
func (n *notifier) Notify() {
	select {
	case n.ch <- struct{}{}:
	default:
	}
}
