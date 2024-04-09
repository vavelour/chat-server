package repos

import (
	"errors"
	"github.com/vavelour/chat/internal/repository/inmemorydb/model"
	"github.com/vavelour/chat/internal/repository/inmemorydb/model/constant"
	"github.com/vavelour/chat/internal/service/mapper"
	"sort"
	"sync"

	"github.com/vavelour/chat/pkg/pagination"

	"github.com/vavelour/chat/internal/domain/entities"
)

var (
	ErrUserIsNotExists = errors.New("this user is not exist")
	ErrChatIsNotExists = errors.New("no chat with this user")
	ErrNonUsers        = errors.New("no users who have written to you")
)

type PrivateDatabase interface {
	Insert(key string, data interface{})
	Get(key string) interface{}
}

type PrivateRepos struct {
	mu sync.RWMutex
	db PrivateDatabase
}

func NewPrivateRepos(db PrivateDatabase) *PrivateRepos {
	return &PrivateRepos{db: db}
}

func (p *PrivateRepos) InsertMessage(m entities.Message) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	members := mapper.MessageToMembersPrivateChat(m)

	data := p.db.Get(constant.UsersKey)

	users, ok := data.(model.UsersTable)
	if !ok {
		return errIncorrectType
	}

	_, ok = users.Table[m.Recipient]
	if !ok {
		return ErrUserIsNotExists
	}

	data = p.db.Get(constant.PrivateChatKey)

	privateChats, ok := data.(model.PrivateChatTable)
	if !ok {
		return errIncorrectType
	}

	chat := privateChats.Table[members]
	chat.Messages = append(chat.Messages, m)
	privateChats.Table[members] = chat

	p.db.Insert(constant.PrivateChatKey, privateChats)

	return nil
}

func (p *PrivateRepos) GetMessages(sender, recipient string, limit, offset int) ([]entities.Message, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	members := mapper.StringToMembersPrivateChat(sender, recipient)

	data := p.db.Get(constant.PrivateChatKey)

	privateChats, ok := data.(model.PrivateChatTable)
	if !ok {
		return nil, errIncorrectType
	}

	messages, ok := privateChats.Table[members]
	if !ok {
		return nil, ErrChatIsNotExists
	}

	paginationMessages, err := pagination.Pagination(messages.Messages, limit, offset)
	if err != nil {
		return nil, err
	}

	return paginationMessages, nil
}

func (p *PrivateRepos) GetUsers(user string) ([]string, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	data := p.db.Get(constant.PrivateChatKey)

	privateChats, ok := data.(model.PrivateChatTable)
	if !ok {
		return nil, errIncorrectType
	}

	var userList model.UserListModel

	for members := range privateChats.Table {
		if members.User1 == user {
			userList.Usernames = append(userList.Usernames, members.User2)
		}

		if members.User2 == user {
			userList.Usernames = append(userList.Usernames, members.User1)
		}
	}

	if len(userList.Usernames) < 1 {
		return nil, ErrNonUsers
	}

	sort.Strings(userList.Usernames)

	return userList.Usernames, nil
}
