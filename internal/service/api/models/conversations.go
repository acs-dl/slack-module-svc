package models

import (
	"github.com/acs-dl/slack-module-svc/internal/data"
	"github.com/acs-dl/slack-module-svc/resources"
)

func NewChatModel(chat data.Conversation) resources.Conversation {
	result := resources.Conversation{
		Key: resources.Key{
			ID:   chat.Id,
			Type: resources.CONVERSATIONS,
		},
		Attributes: resources.ConversationAttributes{
			Title:         chat.Title,
			Id:            chat.Id,
			MembersAmount: chat.MembersAmount,
		},
	}

	return result
}

func NewChatResponse(chat data.Conversation) resources.ConversationResponse {
	return resources.ConversationResponse{
		Data: NewChatModel(chat),
	}
}

func NewChatListModel(chats []data.Conversation) []resources.Conversation {
	newChats := make([]resources.Conversation, 0)

	for _, chat := range chats {
		newChats = append(newChats, NewChatModel(chat))
	}

	return newChats
}
