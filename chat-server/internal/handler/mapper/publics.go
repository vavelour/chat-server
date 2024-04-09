package mapper

import (
	"github.com/vavelour/chat/internal/domain/entities"
	"github.com/vavelour/chat/internal/handler/request"
	"github.com/vavelour/chat/internal/handler/response"
)

func PublicMessageEntitiesToResponse(resp string, messages []entities.Message) response.ShowPublicMessageResponse {
	var res response.ShowPublicMessageResponse
	res.Response = resp

	for _, content := range messages {
		res.Messages = append(res.Messages, content.Content)
	}

	return res
}

func SendPublicMessageRequestToEntities(req request.SendPublicMessageRequest) (m entities.Message) {
	return entities.Message{Sender: req.Sender, Recipient: req.Recipient, Content: req.Content}
}
