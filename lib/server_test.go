package lib

import (
	"net/url"
	"testing"
	"time"

	"github.com/yamavol/greload/test/harness"
)

func Test_negativeDelay(t *testing.T) {
	testDelay := time.Duration(1) * time.Millisecond
	opt := &ServerOptions{
		Host:  &url.URL{Host: "example.com"},
		Port:  12345,
		Delay: testDelay,
	}
	srv := NewServer(opt)
	println(srv.adjustedDelayTime())
	harness.IsTrue(t, testDelay <= defaultDebounceDuration, "delayMs is smaller than default")
	harness.IsTrue(t, srv.adjustedDelayTime() == 0, "adjusted delay time is 0")
}
