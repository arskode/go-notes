package config

import (
	"flag"
	"os"
)

type Config struct {
	DbURL string
}

func NewConfig() *Config {
	conf := &Config{}

	flag.StringVar(&conf.DbURL, "db_url", os.Getenv("DB_URL"), "Database url")
	flag.Parse()

	return conf
}
