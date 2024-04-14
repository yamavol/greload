package lib

import (
	"testing"

	"github.com/yamavol/greload/test/harness"
)

func Test_hasScheme(t *testing.T) {
	harness.IsTrue(t, hasScheme("a://"), "")
	harness.IsTrue(t, hasScheme("http://"), "")
	harness.IsTrue(t, hasScheme("https://"), "")
	harness.IsTrue(t, hasScheme(" h123://"), "")
	harness.IsTrue(t, hasScheme(" chrome-extensions://"), "")
	harness.IsTrue(t, hasScheme(" scheme.with.period://"), "")
	harness.IsTrue(t, hasScheme(" shheme.with.period+plus0://"), "")

	harness.IsFalse(t, hasScheme(" ://"), "")
	harness.IsFalse(t, hasScheme("localhost"), "localhost does not have scheme")
	harness.IsFalse(t, hasScheme("127.0.0.1"), "")
	harness.IsFalse(t, hasScheme("127.0.0.1:8080"), "")
	harness.IsFalse(t, hasScheme("127.0.0.1:3000"), "")
	harness.IsFalse(t, hasScheme("localhost:4000"), "")
	harness.IsFalse(t, hasScheme(":4000"), "")
}

func Test_ParseHost(t *testing.T) {
	u, _ := parseHost("localhost")
	harness.IsEqual(t, u.Scheme, "http", "")
	harness.IsEqual(t, u.Host, "localhost:80", "")
	harness.IsEqual(t, u.Hostname(), "localhost", "")
	harness.IsEqual(t, u.Port(), "80", "")

	u, _ = parseHost("localhost:4000")
	harness.IsEqual(t, u.Scheme, "http", "")
	harness.IsEqual(t, u.Host, "localhost:4000", "")

	u, _ = parseHost(":4000")
	harness.IsEqual(t, u.Scheme, "http", "")
	harness.IsEqual(t, u.Host, "localhost:4000", "")

	u, _ = parseHost("https://12.34.56.78:9012")
	harness.IsEqual(t, u.Scheme, "https", "")
	harness.IsEqual(t, u.Host, "12.34.56.78:9012", "")
}

// func Test_initOptions(t *testing.T) {
// 	opt, err := initServerOptions("localhost:4000", 1234)
// 	harness.IsEqual(t, opt.Port, 1234, "")
// 	harness.IsEqual(t, opt.Forward.String(), "http://localhost:4000", "")

// 	opt, err = initServerOptions("localhost:1000000", -1)
// 	harness.IsEqual(t, opt.Port, 9999, "if port is -1 then use default port")
// 	harness.IsEqual(t, opt.Forward.String(), "http://localhost:1000000", "")

// 	opt, err = initServerOptions("", -1)
// 	harness.IsNotNil(t, err, "if host is empty, return error")
// 	harness.IsEqual(t, opt.Port, -1, "")
// 	harness.IsEqual(t, opt.Forward.String(), "", "")

// }
