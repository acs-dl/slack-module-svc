package models

import (
	"github.com/acs-dl/slack-module-svc/internal/data"
	"github.com/acs-dl/slack-module-svc/resources"
)

func newLink(link string, isExists bool, conversations []data.Conversation) resources.Link {
	return resources.Link{
		Key: resources.Key{
			ID:   link,
			Type: resources.LINKS,
		},
		Attributes: resources.LinkAttributes{
			Link:       link,
			IsExists:   isExists,
			Submodules: NewConversationListModel(conversations),
		},
	}
}

func NewLinkResponse(link string, isExists bool, conversations []data.Conversation) resources.LinkResponse {
	return resources.LinkResponse{
		Data: newLink(link, isExists, conversations),
	}
}
