package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/yamavol/go-argp"
	"github.com/yamavol/greload/lib"
	"github.com/yamavol/greload/log"
)

const (
	flagPort    = "port"
	flagVerbose = "verbose"
	flagLevel   = "level"
	flagWatch   = "watch"
	flagExclude = "exclude"
	flagHelp    = "help"
	flagVersion = "version"

	Version = "0.1.2"
)

var Options = []argp.Option{
	{Short: 'p', Long: flagPort, ArgName: "<port>", Doc: "greload port (http + websocket)"},
	{Short: 'w', Long: flagWatch, ArgName: "<path>", Doc: "add path to watch list"},
	{Short: 'x', Long: flagExclude, ArgName: "<path>", Doc: "add path to ignore list"},
	{Short: 'v', Long: flagVerbose, Flags: argp.OPTION_HIDDEN, Doc: "enable verbose mode"},
	{Short: 'h', Long: flagHelp, Flags: argp.OPTION_HIDDEN, Doc: "print help and exit"},
	{Short: 'V', Long: flagVersion, Flags: argp.OPTION_HIDDEN, Doc: "print version and exit"},
}

func Run() {

	host := ""
	port := lib.DefaultPort
	watch := []string{}
	exclude := []string{}

	result, err := argp.Parse(Options)

	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	if result.HasOpt(flagVerbose) {
		log.SetLogLevel(log.LevelDebug)
	}

	if result.HasOpt(flagHelp) {
		printHelp()
		return
	}

	if result.HasOpt(flagVersion) {
		fmt.Printf("greload %s\n", Version)
		return
	}

	// argument 1 : forwarding host
	if len(result.Args) == 0 {
		printHelp()
		return
	} else {
		host = result.Args[0]
	}

	if result.HasOpt(flagPort) {
		p, err := strconv.Atoi(result.GetOpt(flagPort).Optarg)
		if err != nil {
			log.Errorf("invalid port: %s\n", err)
			return
		}
		if p < 0 || p > 65535 {
			log.Errorf("invalid port: %v\n", p)
			return
		}
		port = p
	}
	for _, opt := range result.Options {
		switch opt.Long {
		case flagWatch:
			watch = append(watch, opt.Optarg)
		case flagExclude:
			exclude = append(exclude, opt.Optarg)
		default:
		}
	}

	// ==============================
	// server options
	// ==============================
	serverOptions := lib.NewServerOption()
	if err = serverOptions.SetForwardHost(host); err != nil {
		log.Error(err, "\n")
		return
	}

	err = serverOptions.SetPort(port)
	if err != nil {
		log.Error(err)
		return
	}

	// ==============================
	// watch options
	// ==============================
	if len(watch) == 0 {
		watch = append(watch, ".")
	}

	proxyServer := lib.NewServer(serverOptions)

	go lib.WatchStart(watch, exclude, proxyServer)

	proxyServer.Start()
}

func printHelp() {
	argp.PrintUsage(os.Stdout, Options, filepath.Base(os.Args[0]), "HOST:PORT")
}
