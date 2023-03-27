package service

import (
	"github.com/axli-personal/drive/backend/drive/adapters"
	"github.com/axli-personal/drive/backend/drive/usecases"
	"github.com/sirupsen/logrus"
	"time"
)

type Config struct {
	MysqlConnectionString string `env:"MYSQL_CONNECTION_STRING" yaml:"mysql-connection-string"`
	RedisConnectionString string `env:"REDIS_CONNECTION_STRING" yaml:"redis-connection-string"`
	UserServiceAddress    string `env:"USER_SERVICE_ADDRESS" yaml:"user-service-address"`
	LogLevel              string `env:"LOG_LEVEL" yaml:"log-level"`
}

type Service struct {
	CreateDrive  usecases.CreateDriveHandler
	CreateFolder usecases.CreateFolderHandler

	StartUpload   usecases.StartUploadHandler
	FinishUpload  usecases.FinishUploadHandler
	StartDownload usecases.StartDownloadHandler

	GetDrive  usecases.GetDriveHandler
	GetFile   usecases.GetFileHandler
	GetFolder usecases.GetFolderHandler
	GetPath   usecases.GetPathHandler

	ShareFile   usecases.ShareFileHandler
	ShareFolder usecases.ShareFolderHandler

	MoveFile   usecases.MoveFileHandler
	MoveFolder usecases.MoveFolderHandler

	RemoveFile   usecases.RemoveFileHandler
	RemoveFolder usecases.RemoveFolderHandler

	DeleteFile usecases.DeleteFileHandler
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

	logger, err := NewLogger(config.LogLevel)
	if err != nil {
		return Service{}, err
	}

	service := Service{
		CreateDrive:  usecases.NewCreateDriveHandler(userService, driveRepo, logger),
		CreateFolder: usecases.NewCreateFolderHandler(userService, driveRepo, folderRepo, logger),

		StartUpload:   usecases.NewStartUploadHandler(userService, driveRepo, fileRepo, logger),
		FinishUpload:  usecases.NewFileUploadedHandler(fileRepo, logger),
		StartDownload: usecases.NewStartDownloadHandler(userService, driveRepo, fileRepo, logger),

		GetDrive:  usecases.NewGetDriveHandler(userService, driveRepo, folderRepo, fileRepo, logger),
		GetFile:   usecases.NewGetFileHandler(userService, driveRepo, fileRepo, logger),
		GetFolder: usecases.NewGetFolderHandler(userService, driveRepo, folderRepo, fileRepo, logger),
		GetPath:   usecases.NewGetPathHandler(userService, driveRepo, folderRepo, logger),

		ShareFile:   usecases.NewShareFileHandler(userService, driveRepo, fileRepo, logger),
		ShareFolder: usecases.NewShareFolderHandler(userService, driveRepo, folderRepo, logger),

		MoveFile:   usecases.NewMoveFileHandler(userService, driveRepo, folderRepo, fileRepo, logger),
		MoveFolder: usecases.NewMoveFolderHandler(userService, driveRepo, folderRepo, logger),

		RemoveFile:   usecases.NewRemoveFileHandler(userService, driveRepo, fileRepo, logger),
		RemoveFolder: usecases.NewRemoveFolderHandler(userService, driveRepo, folderRepo, logger),

		DeleteFile: usecases.NewDeleteFileHandler(userService, driveRepo, fileRepo, eventRepo, logger),
	}

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
