// Package conf module from parse config from env
package conf

import (
	"fmt"

	"github.com/caarlos0/env"
)

type config struct {
	PostgresPassword string `env:"POSTGRES_PASSWORD"`
	PostgresUser     string `env:"POSTGRES_USER"`
	PostgresDB       string `env:"POSTGRES_DB"`
	PostgresHost     string `env:"POSTGRES_HOST"`
	PostgresPort     string `env:"POSTGRES_PORT"`
	RabbitQueueName  string `env:"RABBIT_QUEUE_NAME"`
	RabbitUser       string `env:"RABBIT_USER"`
	RabbitPassword   string `env:"RABBIT_PASSWORD"`
	RabbitHost       string `env:"RABBIT_HOST"`
	RabbitPort       string `env:"RABBIT_PORT"`
}

// NewConfig Get config
func NewConfig() (*config, error) {
	conf := config{}
	err := env.Parse(&conf)
	if err != nil {
		return nil, fmt.Errorf("conf.GetConfig: %v", err)
	}
	return &conf, nil
}
