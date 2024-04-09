package model

import "github.com/vavelour/chat/internal/domain/entities"

type PublicChat struct {
	Messages []entities.Message
}
