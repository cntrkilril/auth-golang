package app

import (
	"github.com/go-playground/validator/v10"
	"time"
)

type (
	Config struct {
		Logger Logger `validate:"required"`
		HTTP   HTTP   `validate:"required"`
		Mongo  Mongo  `validate:"required"`
		Hasher Hasher `validate:"required"`
		Tokens Tokens `validate:"required"`
	}

	Logger struct {
		Level *int8 `validate:"required"`
	}

	HTTP struct {
		Host string `validate:"required"`
		Port string `validate:"required"`
	}

	Mongo struct {
		ConnString string `validate:"required"`
	}

	Hasher struct {
		Cost int `validate:"required"`
	}

	Tokens struct {
		JWTKey                string        `validate:"required"`
		ExpiresInAccessToken  time.Duration `validate:"required"`
		ExpiresInRefreshToken time.Duration `validate:"required"`
	}
)

func LoadConfig() (*Config, error) {

	defaultLogLevel := int8(-1)

	cfg := &Config{
		HTTP: HTTP{
			Host: "localhost",
			Port: "8080",
		},
		Logger: Logger{
			Level: &defaultLogLevel,
		},
		Mongo: Mongo{
			ConnString: "mongodb://root:pass@127.0.0.1:27017",
		},
		Hasher: Hasher{
			Cost: 10,
		},
		Tokens: Tokens{
			JWTKey:                "KZ1NXvjYmYs6oybjaDsqog==",
			ExpiresInRefreshToken: time.Hour,
			ExpiresInAccessToken:  time.Hour,
		},
	}

	err := validator.New().Struct(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
