package model

import "github.com/vavelour/chat/internal/domain/entities"

type UsersTable struct {
	Table map[string]entities.User
}
