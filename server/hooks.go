package main

import (
	"fmt"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

const MessageFormat = `Channel ~%s has been created by %s.`

func (p *Plugin) ChannelHasBeenCreated(c *plugin.Context, channel *model.Channel) {
	p.API.LogInfo("ChannelHasBennCreated", "details", fmt.Sprintf("%#v", channel), "context", fmt.Sprintf("%#v", c))
	if channel.Type != model.CHANNEL_OPEN {
		return
	}

	u, appErr := p.API.GetUser(channel.CreatorId)
	if appErr != nil {
		p.API.LogError("Failed to get user", "details", appErr)
		return
	}
	townSquare, appErr := p.API.GetChannelByName(channel.TeamId, model.DEFAULT_CHANNEL, false)
	if appErr != nil {
		p.API.LogError("Failed to get channel", "details", appErr)
		return
	}

	post := &model.Post{
		Type:      model.POST_DEFAULT,
		ChannelId: townSquare.Id,
		UserId:    p.botUserID,
		Message:   fmt.Sprintf(MessageFormat, channel.Name, u.GetDisplayName(model.SHOW_USERNAME)),
	}

	if _, appErr := p.API.CreatePost(post); appErr != nil {
		p.API.LogError("Failed to create post", "details", appErr)
	}
}
