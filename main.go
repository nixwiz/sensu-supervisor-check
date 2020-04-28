package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/abrander/go-supervisord"
	"github.com/sensu-community/sensu-plugin-sdk/sensu"
	"github.com/sensu/sensu-go/types"
)

// Config represents the check plugin config.
type Config struct {
	sensu.PluginConfig
	Host     string
	Port     int
	Socket   string
	Critical string
}

var (
	plugin = Config{
		PluginConfig: sensu.PluginConfig{
			Name:     "sensu-supervisor-check",
			Short:    "Supervisor Check for Sensu Go",
			Keyspace: "sensu.io/plugins/sensu-supervisor-check/config",
		},
	}

	options = []*sensu.PluginConfigOption{
		{
			Path:      "host",
			Env:       "SUPERVISOR_HOST",
			Argument:  "host",
			Shorthand: "H",
			Default:   "localhost",
			Usage:     "Host running Supervisor",
			Value:     &plugin.Host,
		},
		{
			Path:      "port",
			Env:       "SUPERVISOR_PORT",
			Argument:  "port",
			Shorthand: "P",
			Default:   9001,
			Usage:     "Supervisor listening port",
			Value:     &plugin.Port,
		},
		{
			Path:      "socket",
			Env:       "SUPERVISOR_SOCKET",
			Argument:  "socket",
			Shorthand: "s",
			Default:   "",
			Usage:     "Supervisor listening UNIX domain socket",
			Value:     &plugin.Socket,
		},
		{
			Path:      "critical",
			Env:       "SUPERVISOR_CRITICAL",
			Argument:  "critical",
			Shorthand: "c",
			Default:   "FATAL",
			Usage:     "Supervisor states to consider critical",
			Value:     &plugin.Critical,
		},
	}
)

func main() {
	check := sensu.NewGoCheck(&plugin.PluginConfig, options, checkArgs, executeCheck, false)
	check.Execute()
}

func checkArgs(event *types.Event) (int, error) {
	if len(plugin.Socket) > 0 {
		fi, err := os.Stat(plugin.Socket)
		if err != nil {
			return sensu.CheckStateWarning, err
		}
		if fi.Mode()&os.ModeSocket == 0 {
			return sensu.CheckStateWarning, fmt.Errorf("%s is not a UNIX domain socket", plugin.Socket)
		}

	}
	return sensu.CheckStateOK, nil
}

func executeCheck(event *types.Event) (int, error) {
	var (
		sc  *supervisord.Client
		err error
	)

	if len(plugin.Socket) > 0 {
		sc, err = supervisord.NewUnixSocketClient(plugin.Socket)
		if err != nil {
			fmt.Printf("%s CRITICAL: failed to create client for UNIX domain socket %s", plugin.PluginConfig.Name, plugin.Socket)
			return sensu.CheckStateCritical, nil
		}
	} else {
		url := fmt.Sprintf("http://%s:%d/RPC2", plugin.Host, plugin.Port)
		sc, err = supervisord.NewClient(url)
		if err != nil {
			fmt.Printf("%s CRITICAL: failed to create client for %s", plugin.PluginConfig.Name, url)
			return sensu.CheckStateCritical, nil
		}
	}

	pi, err := sc.GetAllProcessInfo()
	if err != nil {
		fmt.Printf("%s CRITICAL: %v", plugin.PluginConfig.Name, err)
		return sensu.CheckStateCritical, nil
	}

	criticals := strings.Split(strings.ToUpper(plugin.Critical), ",")

	criticalsFound := []string{}
	for _, piEntry := range pi {
		if contains(criticals, piEntry.StateName) {
			criticalsFound = append(criticalsFound, piEntry.Name)
		}
	}

	if len(criticalsFound) > 0 {
		fmt.Printf("%s CRITICAL: process(es) %s not running\n", plugin.PluginConfig.Name, strings.Join(criticalsFound, ", "))
		return sensu.CheckStateCritical, nil
	}

	fmt.Printf("%s OK: All processes running\n", plugin.PluginConfig.Name)
	return sensu.CheckStateOK, nil
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}
