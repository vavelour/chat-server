package mapper

import (
	"github.com/vavelour/chat/internal/domain/entities"
	"github.com/vavelour/chat/internal/repository/postgres/models"
)

func MessageModelToEntities(model models.ChatModel) []entities.Message {
	message := make([]entities.Message, 0)

	for _, val := range model.Messages {
		message = append(message, entities.Message{Sender: val.Sender, Recipient: val.Recipient, Content: val.Content})
	}

	return message
}
