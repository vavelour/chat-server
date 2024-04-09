package repos

import (
	"errors"
	"github.com/vavelour/chat/internal/domain/entities"
	"github.com/vavelour/chat/internal/repository/inmemorydb/model"
	"github.com/vavelour/chat/internal/repository/inmemorydb/model/constant"
	"sync"
)

var (
	errUserAlreadyExists = errors.New("user already exists")
	errIncorrectType     = errors.New("type conversion failed")
	errUnregisteredUser  = errors.New("unregistered user")
)

//go:generate mockgen -source=auth.go -destination=mocks/inmemory_db_mock.go -mock_names=AuthDatabase=MockMemoryDB

type AuthDatabase interface {
	Insert(key string, data interface{})
	Get(key string) interface{}
}

type AuthRepos struct {
	mu sync.RWMutex
	db AuthDatabase
}

func NewAuthRepos(db AuthDatabase) *AuthRepos {
	return &AuthRepos{db: db}
}

func (r *AuthRepos) InsertUser(username, password string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	data := r.db.Get(constant.UsersKey)

	users, ok := data.(model.UsersTable)
	if !ok {
		return errIncorrectType
	}

	_, ok = users.Table[username]
	if ok {
		return errUserAlreadyExists
	}

	users.Table[username] = entities.User{Username: username, Password: password}
	r.db.Insert(constant.UsersKey, users)

	return nil
}

func (r *AuthRepos) GetUser(username string) (entities.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	data := r.db.Get(constant.UsersKey)

	users, ok := data.(model.UsersTable)
	if !ok {
		return entities.User{}, errIncorrectType
	}

	user, ok := users.Table[username]
	if !ok {
		return entities.User{}, errUnregisteredUser
	}

	return user, nil
}
