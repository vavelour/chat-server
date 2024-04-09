package model

import "github.com/vavelour/chat/internal/domain/entities"

type PrivateChat struct {
	Messages []entities.Message
}
