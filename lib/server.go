package lib

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/mattn/go-ieproxy"
	"github.com/yamavol/greload/log"

	"golang.org/x/net/websocket"
)

// ============================================================
// Server Options
// ============================================================
type ServerOptions struct {
	Port int
	Host *url.URL
}

var hasSchemeRe = regexp.MustCompile(`^\s*[0-9A-Za-z.\-\+]+://`)

const DefaultPort int = 9999

func NewServerOption() *ServerOptions {
	return &ServerOptions{
		Port: DefaultPort,
		Host: &url.URL{},
	}
}

func (s *ServerOptions) SetForwardHost(h string) error {
	u, err := parseHost(h)
	if err != nil {
		return err
	}

	s.Host = u
	return nil
}

func (s *ServerOptions) SetPort(port int) error {
	if port < 0 || port > 65535 {
		return fmt.Errorf("port out of range: %v", port)
	}

	s.Port = port
	return nil
}

func hasScheme(s string) bool {
	return hasSchemeRe.Match([]byte(s))
}

func knownPort(scheme string) string {
	switch scheme {
	case "http":
		return "80"
	case "https":
		return "443"
	default:
		panic("unsupported scheme")
	}
}

func parseHost(h string) (*url.URL, error) {
	if len(strings.Trim(h, " \t")) == 0 {
		return nil, errors.New("host undefined")
	}
	if !hasScheme(h) {
		h = "http://" + h
	}
	u, err := url.Parse(h)
	if err != nil {
		return nil, err
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, errors.New("unsupported scheme: " + u.Scheme)
	}
	hostname := u.Hostname()
	port := u.Port()
	if hostname == "" && port == "" {
		return nil, errors.New("host undefined")
	}
	if hostname == "" {
		hostname = "localhost"
	}
	if port == "" {
		port = knownPort(u.Scheme)
	}
	u.Host = hostname + ":" + port
	return u, nil
}

// ============================================================
// HTTP & Websocket server Methods
// ============================================================

// ProxyServer is the libary's main server that handles both HTTP proxying
// and WebSocket connections.
type ProxyServer struct {
	options     ServerOptions
	connections map[*websocket.Conn]struct{}
	lock        sync.Mutex
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

	go func() {
		<-interrupt
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		server.Shutdown(ctx)
	}()

	go func() {
		for range srv.reloadReq.ch {
			srv.handleReload()
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

	rp.ModifyResponse = responseModifier(srv.options.Port)
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

	ws.lock.Lock()
	ws.connections[conn] = struct{}{}
	ws.lock.Unlock()

	log.Debug("[ws] Client connected")

	var msg string
	for {
		err := websocket.Message.Receive(conn, &msg)
		if err != nil {
			ws.lock.Lock()
			delete(ws.connections, conn)
			ws.lock.Unlock()
			log.Debug("[ws] Client disconnected")
			break
		}
	}
}

func (srv *ProxyServer) handleReload() {
	srv.lock.Lock()
	for conn := range srv.connections {
		err := websocket.Message.Send(conn, "reload")
		if err != nil {
			delete(srv.connections, conn)
			log.Errorf("Error broadcasting message to client: %v", err)
		}
	}
	srv.lock.Unlock()
}

// ======================================================================
// notifier (Private)
// ======================================================================
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
