package env

import "flag"

var (
	appname string
	nodeid  string
)

func init() {
	flag.StringVar(&appname, "appname", "", "AppName of application. e.g. -appname=nekoq")
	flag.StringVar(&nodeid, "node", "", "Unique Node Id of application. e.g. -node=nekoq001")
	flag.Parse()
}
