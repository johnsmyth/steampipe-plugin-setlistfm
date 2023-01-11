package main

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/johnsmyth/steampipe-plugin-setlistfm/setlistfm"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{PluginFunc: setlistfm.Plugin})
}
