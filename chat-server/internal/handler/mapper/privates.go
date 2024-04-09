package mapper

import (
	"github.com/vavelour/chat/internal/domain/entities"
	"github.com/vavelour/chat/internal/handler/request"
	"github.com/vavelour/chat/internal/handler/response"
)

func PrivateMessageEntitiesToResponse(resp string, messages []entities.Message) response.ShowPrivateMessageResponse {
	var res response.ShowPrivateMessageResponse
	res.Response = resp

	for _, content := range messages {
		res.Messages = append(res.Messages, content.Content)
	}

	return res
}

func UserListEntitiesToResponse(resp string, users []string) response.ViewUserListResponse {
	var res response.ViewUserListResponse
	res.Response = resp
	res.Messages = append(res.Messages, users...)

	return res
}

func SendPrivateMessageRequestToEntities(req request.SendPrivateMessageRequest) entities.Message {
	return entities.Message{Sender: req.Sender, Recipient: req.Recipient, Content: req.Content}
}
