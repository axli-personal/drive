package service

import (
	"github.com/axli-personal/drive/backend/user/adapters"
	"github.com/axli-personal/drive/backend/user/usecases"
	"github.com/sirupsen/logrus"
	"time"
)

type Config struct {
	MysqlConnectionString string `env:"MYSQL_CONNECTION_STRING" yaml:"mysql-connection-string"`
	RedisConnectionString string `env:"REDIS_CONNECTION_STRING" yaml:"redis-connection-string"`
	LogLevel              string `env:"LOG_LEVEL" yaml:"log-level"`
}

type Service struct {
	Register usecases.RegisterHandler
	Login    usecases.LoginHandler
	GetUser  usecases.GetUserHandler
}

func NewService(config Config) (Service, error) {
	logger, err := NewLogger(config.LogLevel)
	if err != nil {
		return Service{}, err
	}

	userRepo, err := adapters.NewUserRepository(config.MysqlConnectionString)
	if err != nil {
		return Service{}, err
	}

	sessionRepo, err := adapters.NewSessionRepository(config.RedisConnectionString, logger)
	if err != nil {
		return Service{}, err
	}

	return Service{
		Register: usecases.NewRegisterHandler(userRepo, logger),
		Login:    usecases.NewLoginHandler(userRepo, sessionRepo, logger),
		GetUser:  usecases.NewGetUserHandler(userRepo, sessionRepo, logger),
	}, nil
}

func NewLogger(logLevel string) (*logrus.Entry, error) {
	logger := logrus.New()

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return nil, err
	}

	logger.SetLevel(level)

	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.StampMilli,
	})

	return logrus.NewEntry(logger), nil
}
