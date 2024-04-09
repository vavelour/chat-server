package postgres

import (
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/vavelour/chat/pkg/database_utils/postgres"
	"sync"
)

type SqlPostgresDB struct {
	mu sync.RWMutex
	db *sqlx.DB
}

func NewSqlPostgresDB(cfg postgres.SqlPostgresConfig) (*SqlPostgresDB, error) {
	db, err := sqlx.Open("pgx", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &SqlPostgresDB{db: db}, nil
}

func (db *SqlPostgresDB) Get(query string) (*sqlx.Rows, error) {
	rows, err := db.db.Queryx(query)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (db *SqlPostgresDB) Insert(query string) error {
	_, err := db.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
