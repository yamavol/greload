package lib
// ============================================================
// Server Options
// ============================================================


import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)


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
