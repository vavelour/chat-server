package repos

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/vavelour/chat/internal/domain/entities"
	"github.com/vavelour/chat/internal/repository/postgres/mapper"
	"github.com/vavelour/chat/internal/repository/postgres/models"
	"sync"
)

type PrivatePostgresDB interface {
	Insert(query string) error
	Get(query string) (*sqlx.Rows, error)
}

type PrivateSqlRepos struct {
	mu sync.RWMutex
	db PrivatePostgresDB
}

func NewPrivateSqlRepos(db PrivatePostgresDB) *PrivateSqlRepos {
	return &PrivateSqlRepos{db: db}
}

func (p *PrivateSqlRepos) InsertMessage(m entities.Message) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	query := fmt.Sprintf("INSERT INTO private_chats(sender_id, recipient_id, message) "+
		"VALUES ( "+
		"(SELECT id FROM users WHERE username = '%s'), "+
		"(SELECT id FROM users WHERE username = '%s'), "+
		"'%s')", m.Sender, m.Recipient, m.Content)

	if err := p.db.Insert(query); err != nil {
		return err
	}

	return nil
}

func (p *PrivateSqlRepos) GetMessages(sender, recipient string, limit, offset int) ([]entities.Message, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	query := fmt.Sprintf("SELECT sender_id, recipient_id, message "+
		"FROM private_chats "+
		"WHERE (sender_id = (SELECT id FROM users WHERE username = '%s') AND recipient_id = (SELECT id FROM users WHERE username = '%s')) "+
		"OR (sender_id = (SELECT id FROM users WHERE username = '%s') AND recipient_id = (SELECT id FROM users WHERE username = '%s')) "+
		"LIMIT %d OFFSET %d", sender, recipient, recipient, sender, limit, offset)

	rows, err := p.db.Get(query)
	if err != nil {
		return nil, err
	}

	chat := models.ChatModel{Messages: make([]models.MessageModel, 0)}
	for rows.Next() {
		var message models.MessageModel
		err := rows.StructScan(&message)
		if err != nil {
			return nil, err
		}

		chat.Messages = append(chat.Messages, message)
	}

	return mapper.MessageModelToEntities(chat), nil

}

func (p *PrivateSqlRepos) GetUsers(user string) ([]string, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	var userList models.UserListModel
	query := fmt.Sprintf("SELECT DISTINCT u.username "+
		"FROM users u "+
		"JOIN private_chats pc ON u.id = pc.sender_id "+
		"WHERE pc.recipient_id = (SELECT id FROM users WHERE username = '%s') "+
		"UNION "+
		"SELECT DISTINCT u.username "+
		"FROM users u "+
		"JOIN private_chats pc ON u.id = pc.recipient_id "+
		"WHERE pc.sender_id = (SELECT id FROM users WHERE username = '%s')", user, user)

	rows, err := p.db.Get(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var user string
		err := rows.Scan(&user)
		if err != nil {
			return nil, err
		}

		userList.Usernames = append(userList.Usernames, user)
	}

	return userList.Usernames, nil
}
