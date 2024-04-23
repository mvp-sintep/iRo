package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// EventConfig - блок данных конфигурации события
type EventConfig struct {
	Text   string `yaml:"text"` // текст сообщения, используется для поиска
	Active string `yaml:"act"`  // адрес бита активации в ядре данных приложения
	ACK    string `yaml:"ack"`  // адрес бита подтверждения в ядре данных приложения
	SHLVD  string `yaml:"shlvd"`
	DSUPR  string `yaml:"dsupr"`
	OOSRV  string `yaml:"oosrv"`
	Store  int    `yaml:"store"`
}

// AEConfig - настройка модуля тревог и событий
type AEConfig struct {
	Enabled bool
	Event   []EventConfig
}

// NewAEConfig - создание блока данных настройки
func NewAEConfig() *AEConfig {
	return &AEConfig{
		Enabled: false,
		Event:   []EventConfig{},
	}
}

// Load - загрузка данных конфигурации
func (o *AEConfig) Load(path *string) error {
	file, err := os.ReadFile(*path)
	if err == nil {
		err = yaml.Unmarshal(file, o)
	}
	return err
}
