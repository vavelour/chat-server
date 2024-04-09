package configs

import "github.com/spf13/viper"

type Config struct {
	DB     DBConfig
	Server ServerConfig
	Auth   AuthConfig
}

func InitConfig() (Config, error) {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, err
	}

	cfg := Config{
		DB: DBConfig{
			Type:     viper.GetString("db.type"),
			Host:     viper.GetString("db.host"),
			Port:     viper.GetString("db.port"),
			User:     viper.GetString("db.user"),
			DBName:   viper.GetString("db.db_name"),
			Password: viper.GetString("db.password"),
			SSLMode:  viper.GetString("db.ssl_mode"),
		},
		Server: ServerConfig{
			Addr:           viper.GetString("server.addr"),
			ReadTimeout:    viper.GetDuration("server.read_timeout"),
			WriteTimeout:   viper.GetDuration("server.write_timeout"),
			MaxHeaderBytes: viper.GetInt("server.max_header_bytes"),
		},
		Auth: AuthConfig{Type: viper.GetString("auth.type")},
	}

	return cfg, nil
}
