package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Configuration struct {
	Server struct {
		Addr string `default:"localhost"`
		Port string `default:"5000"`
	}

	DB struct {
		Host     string `default:"localhost"`
		Port     int    `default:"5433"`
		Name     string `default:"findhotel"`
		User     string `default:"findhotel"`
		Password string `default:"findhotel"`
		SSL      string `default:"disable"`
	}

	Importer struct {
		Filepath string `default:"./data_dump.csv"`
	}

	LogLevel  string `default:"debug"`
	PrettyLog bool   `default:"true"`
}

var Config Configuration

const prefix = "findhotel"

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
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSL)
}
