package config

import "github.com/go-playground/validator/v10"

type Config struct {
	DB     DB     `yaml:"db"`
	Logger Logger `yaml:"logger"`
	Env    string `yaml:"env"`
}

type DB struct {
	Downloads string `yaml:"downloads" validate:"required"`
	Queues    string `yaml:"queues" validate:"required"`
}

type Logger struct {
	Level string `yaml:"level" validate:"required,oneof=trace debug info warn error fatal"`
}

func (c Config) Validate() error {
	v := validator.New()
	return v.Struct(c)
}
