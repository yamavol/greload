package greload

import "github.com/yamavol/greload/lib"

// greload is a CLI tool, but you can import it as a library if you wish.

// WithReload is a http middleware function that installs greload to your
// webserver
var WithReload = lib.WithReload

// WatchStart is a function to start watching
var WatchStart = lib.WatchStart

// InjectHtml returns a HTML with reload script injected (if enabled)
var InjectHtml = lib.InjectHtml

// WriteHtml writes InjectHTML to the [ResponseWriter].
var WriteHtml = lib.WriteHtml
