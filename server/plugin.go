package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"github.com/pkg/errors"
)

const (
	BotUsername    = "announce-new-channel"
	BotDisplayName = "announce-new-channel"
	BotDescription = "Announce new channel when it's created."
)

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type Plugin struct {
	plugin.MattermostPlugin
	botUserID string

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration
}

func (p *Plugin) OnActivate() error {
	botUserID, appErr := p.Helpers.EnsureBot(&model.Bot{
		Username:    BotUsername,
		DisplayName: BotDisplayName,
		Description: BotDescription,
	}, plugin.ProfileImagePath("assets/announce-new-channel-icon.png"))
	if appErr != nil {
		return errors.Wrap(appErr, "failed to ensure bot user")
	}
	p.botUserID = botUserID

	return nil
}

// ServeHTTP demonstrates a plugin that handles HTTP requests by greeting the world.
func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, world!")
}

// See https://developers.mattermost.com/extend/plugins/server/reference/
