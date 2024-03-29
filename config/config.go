package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Configuration struct {
	Server struct {
		Addr string `default:"localhost"`
		Port string `default:"5000"`

		MaxUploadMB int64 `default:"200"`
	}

	DB struct {
		Host     string `default:"localhost"`
		Port     int    `default:"5432"`
		Name     string `default:"geoservice"`
		User     string `default:"geoservice"`
		Password string `default:"geoservice"`
		SSL      string `default:"disable"`
	}

	Importer struct {
		Filepath string `default:"./data_dump.csv"`
	}

	LogLevel  string `default:"debug"`
	PrettyLog bool   `default:"true"`
}

var Config Configuration

const prefix = "geoservice"

func Init() error {
	return envconfig.Process(prefix, &Config)
}

func Usage() {
	err := envconfig.Usage(prefix, &Config)
	if err != nil {
		panic(err)
	}
}

func DSN() string {
	cfg := Config.DB
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name, cfg.SSL)
}
