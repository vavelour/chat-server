package configs

import "time"

type DBConfig struct {
	Type     string
	Host     string
	Port     string
	User     string
	DBName   string
	Password string
	SSLMode  string
}

type ServerConfig struct {
	Addr           string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	MaxHeaderBytes int
}

type AuthConfig struct {
	Type string
}
