package repos

import (
	"github.com/vavelour/chat/internal/repository/inmemorydb/model"
	"github.com/vavelour/chat/internal/repository/inmemorydb/model/constant"
	"sync"

	"github.com/vavelour/chat/pkg/pagination"

	"github.com/vavelour/chat/internal/domain/entities"
)

type PublicDatabase interface {
	Insert(key string, data interface{})
	Get(key string) interface{}
}

type PublicRepos struct {
	mu sync.RWMutex
	db PublicDatabase
}

func NewPublicRepos(db PublicDatabase) *PublicRepos {
	return &PublicRepos{db: db}
}

func (pub *PublicRepos) InsertMessage(m entities.Message) error {
	pub.mu.Lock()
	defer pub.mu.Unlock()

	data := pub.db.Get(constant.PublicChatKey)

	publicMessages, ok := data.(model.PublicChat)
	if !ok {
		return errIncorrectType
	}

	publicMessages.Messages = append(publicMessages.Messages, m)
	pub.db.Insert(constant.PublicChatKey, publicMessages)

	return nil
}

func (pub *PublicRepos) GetMessages(limit, offset int) ([]entities.Message, error) {
	pub.mu.RLock()
	defer pub.mu.RUnlock()

	data := pub.db.Get(constant.PublicChatKey)

	publicMessages, ok := data.(model.PublicChat)
	if !ok {
		return nil, errIncorrectType
	}

	paginationMessages, err := pagination.Pagination(publicMessages.Messages, limit, offset)
	if err != nil {
		return nil, err
	}

	return paginationMessages, nil
}
