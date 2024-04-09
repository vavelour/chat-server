package mapper

import (
	"github.com/vavelour/chat/internal/repository/inmemorydb/model"
	"sort"

	"github.com/vavelour/chat/internal/domain/entities"
)

func MessageToMembersPrivateChat(m entities.Message) model.MembersPrivateChatModel {
	return StringToMembersPrivateChat(m.Sender, m.Recipient)
}

func StringToMembersPrivateChat(username1, username2 string) model.MembersPrivateChatModel {
	usernames := []string{username1, username2}
	sort.Strings(usernames)

	return model.MembersPrivateChatModel{User1: usernames[0], User2: usernames[1]}
}
