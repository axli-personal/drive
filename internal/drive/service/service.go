package service

import (
	"context"
	"github.com/axli-personal/drive/internal/drive/adapters"
	"github.com/axli-personal/drive/internal/drive/usecases"
	"github.com/sirupsen/logrus"
	"time"
)

type Config struct {
	MysqlConnectionString string   `env:"MYSQL_CONNECTION_STRING" yaml:"mysql-connection-string"`
	RedisConnectionString string   `env:"REDIS_CONNECTION_STRING" yaml:"redis-connection-string"`
	UserServiceAddress    string   `env:"USER_SERVICE_ADDRESS" yaml:"user-service-address"`
	LogLevel              string   `env:"LOG_LEVEL" yaml:"log-level"`
	StorageEndpoints      []string `env:"STORAGE_ENDPOINTS" yaml:"storage-endpoints"`
}

type Service struct {
	CreateDrive    usecases.CreateDriveHandler
	CreateFile     usecases.CreateFileHandler
	CreateFolder   usecases.CreateFolderHandler
	CanUpload      usecases.CanUploadHandler
	CanDownload    usecases.CanDownloadHandler
	FileUploaded   usecases.FileUploadedHandler
	FileDownloaded usecases.FileDownloadedHandler
	GetDrive       usecases.GetDriveHandler
	GetFile        usecases.GetFileHandler
	GetFolder      usecases.GetFolderHandler
	GetPath        usecases.GetPathHandler
	MoveFile       usecases.MoveFileHandler
	MoveFolder     usecases.MoveFolderHandler
	RemoveFile     usecases.RemoveFileHandler
	RemoveFolder   usecases.RemoveFolderHandler
	DeleteFile     usecases.DeleteFileHandler
}

func NewService(config Config) (Service, error) {
	driveRepo, err := adapters.NewMysqlDriveRepository(config.MysqlConnectionString)
	if err != nil {
		return Service{}, err
	}

	folderRepo, err := adapters.NewMysqlFolderRepository(config.MysqlConnectionString)
	if err != nil {
		return Service{}, err
	}

	fileRepo, err := adapters.NewMysqlFileRepository(config.MysqlConnectionString)
	if err != nil {
		return Service{}, err
	}

	eventRepo, err := adapters.NewRedisEventRepository(config.RedisConnectionString)
	if err != nil {
		return Service{}, err
	}

	userService, err := adapters.NewRPCUserService(config.UserServiceAddress)
	if err != nil {
		return Service{}, err
	}

	storageCluster, err := adapters.NewStorageScheduler(config.RedisConnectionString, config.StorageEndpoints)
	if err != nil {
		return Service{}, err
	}

	logger, err := NewLogger(config.LogLevel)
	if err != nil {
		return Service{}, err
	}

	service := Service{
		CreateDrive:    usecases.NewCreateDriveHandler(userService, driveRepo, logger),
		CreateFile:     usecases.NewCreateFileHandler(userService, storageCluster, driveRepo, fileRepo, logger),
		CreateFolder:   usecases.NewCreateFolderHandler(userService, driveRepo, folderRepo, logger),
		CanUpload:      usecases.NewCanUploadHandler(userService, driveRepo, fileRepo, logger),
		CanDownload:    usecases.NewCanDownloadHandler(userService, driveRepo, fileRepo, logger),
		FileUploaded:   usecases.NewFileUploadedHandler(fileRepo, logger),
		FileDownloaded: usecases.NewFileDownloadedHandler(fileRepo, logger),
		GetDrive:       usecases.NewGetDriveHandler(userService, driveRepo, folderRepo, fileRepo, logger),
		GetFile:        usecases.NewGetFileHandler(userService, driveRepo, fileRepo, logger),
		GetFolder:      usecases.NewGetFolderHandler(userService, driveRepo, folderRepo, fileRepo, logger),
		GetPath:        usecases.NewGetPathHandler(userService, driveRepo, folderRepo, logger),
		MoveFile:       usecases.NewMoveFileHandler(userService, driveRepo, folderRepo, fileRepo, logger),
		MoveFolder:     usecases.NewMoveFolderHandler(userService, driveRepo, folderRepo, logger),
		RemoveFile:     usecases.NewRemoveFileHandler(userService, driveRepo, fileRepo, logger),
		RemoveFolder:   usecases.NewRemoveFolderHandler(userService, driveRepo, folderRepo, eventRepo, logger),
		DeleteFile:     usecases.NewDeleteFileHandler(userService, driveRepo, fileRepo, eventRepo, logger),
	}

	go func() {
		for {
			ctx := context.Background()

			event, err := eventRepo.GetFileUploaded(ctx)
			if err != nil {
				logrus.WithError(err).Info("get event from file uploaded stream")
				time.Sleep(50 * time.Millisecond)
				continue
			}

			result, err := service.FileUploaded.Handle(ctx, event)
			if result.CanAcknowledge {
				eventRepo.AckGetFileUploaded(ctx, event.EventId)
			}
		}
	}()

	go func() {
		for {
			ctx := context.Background()

			event, err := eventRepo.GetFileDownloaded(ctx)
			if err != nil {
				logrus.WithError(err).Info("get event from file downloaded stream")
				time.Sleep(50 * time.Millisecond)
				continue
			}

			result, err := service.FileDownloaded.Handle(ctx, event)
			if result.CanAcknowledge {
				eventRepo.AckFileDownloaded(ctx, event.EventId)
			}
		}
	}()

	return service, nil
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
