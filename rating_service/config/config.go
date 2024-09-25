package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Database struct {
		Host     string `yaml:"host"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
		SSLMode  string `yaml:"sslmode"`
	} `yaml:"database"`
	RabbitMQ struct {
		URL       string `yaml:"url"`
		QueueName string `yaml:"queue_name"`
	} `yaml:"rabbitmq"`
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`
}

func LoadConfig(file string) (*Config, error) {
	var cfg Config
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("could not read config file: %w", err)
	}
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal config file: %w", err)
	}
	return &cfg, nil
}

func (c *Config) DSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host, c.Database.User, c.Database.Password, c.Database.DBName, c.Database.SSLMode)
}
