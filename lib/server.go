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
// Server Methods
// ============================================================

// Start a HTTP & WebSocket server
func ServerStart(option *ServerOptions) {

	server := http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%v", option.Port),
		Handler: http.HandlerFunc(serverHandler(option)),
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	go func() {
		<-interrupt
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		server.Shutdown(ctx)
	}()

	log.Info("reload server is running on port", option.Port)
	log.Info("redirecting access to", option.Host.String())
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error("Error:", err)
	}
}

func serverHandler(option *ServerOptions) func(w http.ResponseWriter, r *http.Request) {

	// reverse proxy server config
	rp := &httputil.ReverseProxy{}

	rp.ModifyResponse = responseModifier(option.Port)
	rp.Transport = &http.Transport{
		// use http_proxy (env) if set, otherwise use system proxy
		Proxy:             ieproxy.GetProxyFunc(),
		DisableKeepAlives: true,
	}
	rp.Rewrite = func(pr *httputil.ProxyRequest) {
		pr.SetXForwarded()
		pr.SetURL(option.Host)
	}

	// websocket server config
	ws := websocket.Handler(websockHandler)

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Upgrade") == "websocket" {
			ws.ServeHTTP(w, r)
		} else {
			rp.ServeHTTP(w, r)
		}
	}
}

// ======================================================================
// WebSocket Server related
// ======================================================================

var connections = make(map[*websocket.Conn]struct{})

func websockHandler(ws *websocket.Conn) {
	defer ws.Close()
	connections[ws] = struct{}{}
	log.Debug("[ws] Client connected")

	var msg string
	for {
		err := websocket.Message.Receive(ws, &msg)
		if err != nil {
			delete(connections, ws)
			log.Debug("[ws] Client disconnected")
			break
		}
	}
}

func broadcastReload() {
	for conn := range connections {
		err := websocket.Message.Send(conn, "reload")
		if err != nil {
			delete(connections, conn)
			log.Errorf("Error broadcasting message to client: %v", err)
		}
	}
}
