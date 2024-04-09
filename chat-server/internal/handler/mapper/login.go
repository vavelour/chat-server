package mapper

import (
	"github.com/vavelour/chat/internal/domain/entities"
	"github.com/vavelour/chat/internal/handler/request"
)

func BasicLogInRequestToEntities(req request.BasicAuthLogInRequest) entities.User {
	return entities.User{Username: req.Username, Password: req.Password}
}
