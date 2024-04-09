package repos

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/vavelour/chat/internal/domain/entities"
	"github.com/vavelour/chat/internal/repository/postgres/models"
	"sync"
)

//go:generate mockgen -source=auth.go -destination=mocks/postgres_db_mock.go -mock_names=AuthPostgresDB=MockPostgresDB

type AuthPostgresDB interface {
	Insert(query string) error
	Get(query string) (*sqlx.Rows, error)
}

type AuthSqlRepos struct {
	mu sync.RWMutex
	db AuthPostgresDB
}

func NewAuthSqlRepos(db AuthPostgresDB) *AuthSqlRepos {
	return &AuthSqlRepos{db: db}
}

func (r *AuthSqlRepos) InsertUser(username, password string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	query := fmt.Sprintf("INSERT INTO users(username, password_hash) VALUES('%s', '%s')", username, password)

	if err := r.db.Insert(query); err != nil {
		return err
	}

	return nil
}

func (r *AuthSqlRepos) GetUser(username string) (entities.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	query := fmt.Sprintf("SELECT username, password_hash FROM users WHERE username = '%s'", username)
	var user models.UserModel

	rows, err := r.db.Get(query)
	if err != nil {
		return entities.User{}, err
	}

	for rows.Next() {
		rows.StructScan(&user)
	}

	return entities.User{Username: user.Username, Password: user.Password}, nil
}
