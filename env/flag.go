package env

import (
	"flag"
	"fmt"
	"os"
)

var (
	appname string
	nodeid  string
)

func init() {
	flag.StringVar(&appname, "appname", "", "AppName of application. e.g. -appname=nekoq")
	flag.StringVar(&nodeid, "node", "", "Unique Node Id of application. e.g. -node=nekoq001")
	flag.Parse()
	ensure, found := os.LookupEnv("NEKO_NOT_ENSURE_ENV")
	if found && ensure == "true" {
		return
	} else {
		EnsureEnvFlag()
		return
	}
}

func EnsureEnvFlag() {
	if appname == "" {
		fmt.Fprintln(os.Stderr, "Please set AppName using flag '-appname'.")
		os.Exit(-100)
	}
	if nodeid == "" {
		fmt.Fprintln(os.Stderr, "Please set NodeId using flag '-node'.")
		os.Exit(-101)
	}
}
