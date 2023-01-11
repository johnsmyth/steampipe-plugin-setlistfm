package setlistfm

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/schema"
)

type setlistfmConfig struct {
	Token *string `cty:"token"`
}

var ConfigSchema = map[string]*schema.Attribute{
	"token": {
		Type: schema.TypeString,
	},
}

func ConfigInstance() interface{} {
	return &setlistfmConfig{}
}

// GetConfig :: retrieve and cast connection config from query data
func GetConfig(connection *plugin.Connection) setlistfmConfig {
	if connection == nil || connection.Config == nil {
		return setlistfmConfig{}
	}
	config, _ := connection.Config.(setlistfmConfig)
	return config
}