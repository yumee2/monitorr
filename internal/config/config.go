package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Services []Service `yaml:"services"`
	Telegram Telegram  `yaml:"telegram"`
}

type Service struct {
	Name     string `yaml:"name"`
	URL      string `yaml:"url"`
	Interval int    `yaml:"interval"`
}

type Telegram struct {
	BotToken         string `yaml:"botToken"`
	ChatID           string `yaml:"chatID"`
	FailureThreshold int    `yaml:"failureThreshold"`
}

func ReadConfig(path string) (*Config, error) {
	var cfg Config
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
