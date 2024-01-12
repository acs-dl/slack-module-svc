package models

import (
	"github.com/acs-dl/slack-module-svc/internal/data"
	"github.com/acs-dl/slack-module-svc/resources"
)

func NewConversationModel(conversation data.Conversation) resources.Conversation {
	result := resources.Conversation{
		Key: resources.Key{
			ID:   conversation.Id,
			Type: resources.CONVERSATIONS,
		},
		Attributes: resources.ConversationAttributes{
			Title:         conversation.Title,
			Id:            conversation.Id,
			MembersAmount: conversation.MembersAmount,
		},
	}

	return result
}

func NewConversatioResponse(conversation data.Conversation) resources.ConversationResponse {
	return resources.ConversationResponse{
		Data: NewConversationModel(conversation),
	}
}

func NewConversationListModel(conversations []data.Conversation) []resources.Conversation {
	newConversations := make([]resources.Conversation, 0)

	for _, conversation := range conversations {
		newConversations = append(newConversations, NewConversationModel(conversation))
	}

	return newConversations
}
