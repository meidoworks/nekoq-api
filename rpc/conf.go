package rpc

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"import.moetang.info/go/nekoq-api/errorutil"
)

var confFilePath string
var serviceConfig *ServiceConfig
var enabledService map[string]ServiceItem

func init() {
	flag.StringVar(&confFilePath, "rpcfile", "", "Rpc configuration file path. -rpcfile=services.conf")
	flag.Parse()
	serviceConfig = &ServiceConfig{}
	enabledService = make(map[string]ServiceItem)
	if confFilePath != "" {
		initConfig()
	}
}

type ServiceConfig struct {
	Services []ServiceItem `json:"services"`
}

type ServiceItem struct {
	ServiceName string            `json:"name"`
	Enable      bool              `json:"enable"`
	Method      []string          `json:"methods"`
	Config      map[string]string `json:"config"`
}

func initConfig() {
	f, err := os.Open(confFilePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, errorutil.NewNested("open rpc config file error.", err))
		os.Exit(-111)
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Fprintln(os.Stderr, errorutil.NewNested("read rpc config file error.", err))
		os.Exit(-112)
	}
	err = json.Unmarshal(data, serviceConfig)
	if err != nil {
		fmt.Fprintln(os.Stderr, errorutil.NewNested("unmarshal rpc config file data error.", err))
		os.Exit(-113)
	}
	for _, v := range serviceConfig.Services {
		if v.Enable {
			enabledService[v.ServiceName] = v
		}
	}
}
