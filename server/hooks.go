package main

import (
	"fmt"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

const MessageFormat = `Channel ~%s has been created by @%s.`

func (p *Plugin) ChannelHasBeenCreated(c *plugin.Context, channel *model.Channel) {
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
	}

	switch p.getConfiguration().MessageType {
	case "message-attachments":
		purpose := channel.Purpose
		if len(purpose) == 0 {
			purpose = "(no purpose)"
		}
		header := channel.Header
		if len(header) == 0 {
			header = "(no header)"
		}
		attachments := &model.SlackAttachment{
			Color:      "#FF8000",
			AuthorName: "mattermost-plugin-announce-new-channel",
			AuthorLink: "https://github.com/kaakaa/mattermost-plugin-announce-new-channel",
			Title:      fmt.Sprintf(MessageFormat, channel.Name, u.GetDisplayName(model.SHOW_USERNAME)),
			Fields: []*model.SlackAttachmentField{{
				Title: "Purpose",
				Value: purpose,
				Short: false,
			}, {
				Title: "Header",
				Value: header,
				Short: false,
			}},
		}
		post.AddProp("attachments", []*model.SlackAttachment{attachments})
	case "simple":
		fallthrough
	default:
		post.Message = fmt.Sprintf(MessageFormat, channel.Name, u.GetDisplayName(model.SHOW_USERNAME))
	}

	if _, appErr := p.API.CreatePost(post); appErr != nil {
		p.API.LogError("Failed to create post", "details", appErr)
	}
}
