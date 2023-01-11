package setlistfm

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name: "steampipe-plugin-setlistfm",
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
			Schema:      ConfigSchema,
		},
		// DefaultGetConfig: &plugin.GetConfig{
		// 	ShouldIgnoreError: errors.NotFoundError,
		// },
		DefaultTransform: transform.FromGo(),
		TableMap: map[string]*plugin.Table{
			"setlistfm_setlist": tableSetlistFMSetlist(ctx),
		},
	}
	return p
}
