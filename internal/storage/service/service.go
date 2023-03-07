package service

import (
	"github.com/axli-personal/drive/internal/storage/adapters"
	"github.com/axli-personal/drive/internal/storage/usecases"
	"github.com/sirupsen/logrus"
	"time"
)

type Config struct {
	Endpoint              string `env:"ENDPOINT" yaml:"endpoint"`
	DataDirectory         string `env:"DATA_DIRECTORY" yaml:"data-directory"`
	RequestPerSecond      int    `env:"REQUEST_PER_SECOND" yaml:"request-per-second"`
	DriveServiceAddress   string `env:"DRIVE_SERVICE_ADDRESS" yaml:"drive-service-address"`
	RedisConnectionString string `env:"REDIS_CONNECTION_STRING" yaml:"redis-connection-string"`
	LogLevel              string `env:"LOG_LEVEL" yaml:"log-level"`
}

type Service struct {
	UploadObject   usecases.UploadObjectHandler
	DownloadObject usecases.DownloadObjectHandler
}

func NewService(config Config) (Service, error) {
	driveService, err := adapters.NewRPCDriveService(config.DriveServiceAddress)
	if err != nil {
		return Service{}, err
	}

	objectRepo, err := adapters.NewDiskObjectRepository(config.DataDirectory)
	if err != nil {
		return Service{}, err
	}

	eventRepo, err := adapters.NewRedisEventRepository(config.RedisConnectionString)
	if err != nil {
		return Service{}, err
	}

	capacityRepo, err := adapters.NewRedisCapacityRepository(config.RedisConnectionString, config.Endpoint, config.DataDirectory, config.RequestPerSecond)
	if err != nil {
		return Service{}, err
	}

	logger, err := NewLogger(config.LogLevel)
	if err != nil {
		return Service{}, err
	}

	return Service{
		UploadObject:   usecases.NewUploadObjectHandler(config.Endpoint, driveService, objectRepo, eventRepo, capacityRepo, logger),
		DownloadObject: usecases.NewDownloadObjectHandler(driveService, objectRepo, eventRepo, capacityRepo, logger),
	}, nil
}

func NewLogger(logLevel string) (*logrus.Entry, error) {
	logger := logrus.New()

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return nil, err
	}

	logger.SetLevel(level)

	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.StampMilli,
	})

	return logrus.NewEntry(logger), nil
}
