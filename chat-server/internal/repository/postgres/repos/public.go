package repos

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/vavelour/chat/internal/domain/entities"
	"github.com/vavelour/chat/internal/repository/postgres/mapper"
	"github.com/vavelour/chat/internal/repository/postgres/models"
	"sync"
)

type PublicPostgresDB interface {
	Insert(query string) error
	Get(query string) (*sqlx.Rows, error)
}

type PublicSqlRepos struct {
	mu sync.RWMutex
	db PublicPostgresDB
}

func NewPublicSqlRepos(db PublicPostgresDB) *PublicSqlRepos {
	return &PublicSqlRepos{db: db}
}

func (pub *PublicSqlRepos) InsertMessage(m entities.Message) error {
	pub.mu.Lock()
	defer pub.mu.Unlock()

	query := fmt.Sprintf("INSERT INTO global_chat(sender_id, message) VALUES ((SELECT id FROM users WHERE username = '%s'), '%s')", m.Sender, m.Content)

	if err := pub.db.Insert(query); err != nil {
		return err
	}

	return nil
}

func (pub *PublicSqlRepos) GetMessages(limit, offset int) ([]entities.Message, error) {
	pub.mu.RLock()
	defer pub.mu.RUnlock()

	query := fmt.Sprintf("SELECT sender_id, message FROM global_chat LIMIT %d OFFSET %d", limit, offset)

	rows, err := pub.db.Get(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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
