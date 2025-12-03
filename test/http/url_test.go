package http_test

import (
	"net/url"
	"testing"

	"github.com/yamavol/greload/test/harness"
)

func Test_UrlParse(t *testing.T) {
	u, err := url.Parse("abcdefg://localhost")
	if err != nil {
		println(err)
	}
	harness.IsNil(t, err, "parse url should success")
	harness.IsEqual(t, u.Scheme, "abcdefg", "")
	harness.IsEqual(t, u.Hostname(), "localhost", "")
	harness.IsEqual(t, u.Port(), "", "")
}

func Test_UrlParseWithUserinfo(t *testing.T) {
	u, err := url.Parse("abc-def://user:p%40ss@proxy:3128")

	harness.IsNil(t, err, "parse url should success")
	harness.IsEqual(t, u.Scheme, "abc-def", "")
	harness.IsEqual(t, u.Host, "proxy:3128", "")
	harness.IsEqual(t, u.User.String(), "user:p%40ss", "")
	harness.IsEqual(t, u.User.Username(), "user", "")

	pass, exist := u.User.Password()
	harness.IsTrue(t, exist, "")
	harness.IsEqual(t, pass, "p@ss", "")
}
