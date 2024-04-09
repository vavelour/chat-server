package postgres

type SqlPostgresConfig struct {
	Host     string
	Port     string
	User     string
	DBName   string
	Password string
	SSLMode  string
}
