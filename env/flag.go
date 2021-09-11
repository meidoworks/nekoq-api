package env

import (
	"flag"
	"fmt"
	"os"
	"time"
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
		_, _ = fmt.Fprintln(os.Stderr, "Please set AppName using flag '-appname'.")
		n := time.Now()
		appname = fmt.Sprintf("sampleapp_%d", n.Unix())
	}
	if nodeid == "" {
		_, _ = fmt.Fprintln(os.Stderr, "Please set NodeId using flag '-node'.")
		n := time.Now()
		appname = fmt.Sprintf("samplenode_%d", n.Unix())
	}
}
