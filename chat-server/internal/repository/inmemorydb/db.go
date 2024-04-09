package inmemorydb

import (
	"github.com/vavelour/chat/internal/repository/inmemorydb/model"
	"github.com/vavelour/chat/internal/repository/inmemorydb/model/constant"
	"sync"

	"github.com/vavelour/chat/internal/domain/entities"
)

type MemoryDB struct {
	mu sync.RWMutex
	db map[string]interface{}
}

func NewDB() *MemoryDB {
	db := make(map[string]interface{})

	db[constant.UsersKey] = model.UsersTable{Table: make(map[string]entities.User)}
	db[constant.PublicChatKey] = model.PublicChat{Messages: make([]entities.Message, 0)}
	db[constant.PrivateChatKey] = model.PrivateChatTable{Table: make(map[model.MembersPrivateChatModel]model.PrivateChat)}

	return &MemoryDB{db: db}
}

func (db *MemoryDB) Insert(key string, data interface{}) {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.db[key] = data
}

func (db *MemoryDB) Get(key string) interface{} {
	db.mu.RLock()
	defer db.mu.RUnlock()

	return db.db[key]
}
