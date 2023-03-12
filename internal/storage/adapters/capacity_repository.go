package adapters

import (
	"context"
	"errors"
	"github.com/axli-personal/drive/internal/storage/repository"
	"github.com/redis/go-redis/v9"
	"syscall"
	"time"
)

var (
	ErrNoRequestCapacity = errors.New("no request capacity")
)

const (
	KeyDiskCapacity    = "DiskCapacity"
	KeyRequestCapacity = "RequestCapacity"
)

type RedisCapacityRepository struct {
	client           *redis.Client
	endpoint         string
	directoryPath    string
	requestPerSecond int
}

func NewRedisCapacityRepository(connectionString string, endpoint string, directoryPath string, requestPerSecond int) (repository.CapacityRepository, error) {
	options, err := redis.ParseURL(connectionString)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(options)

	err = client.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}

	repo := RedisCapacityRepository{
		client:           client,
		endpoint:         endpoint,
		directoryPath:    directoryPath,
		requestPerSecond: requestPerSecond,
	}

	go func() {
		requestCapacityTicker := time.NewTicker(1 * time.Second)

		for range requestCapacityTicker.C {
			repo.updateRequestCapacity(context.Background())
		}
	}()

	go func() {
		diskCapacityTicker := time.NewTicker(5 * time.Second)

		for range diskCapacityTicker.C {
			repo.updateDiskCapacity(context.Background())
		}
	}()

	return repo, nil
}

func (repo RedisCapacityRepository) DecreaseRequestCapacity(ctx context.Context) error {
	capacity, err := repo.client.Get(ctx, repo.endpoint+":"+KeyRequestCapacity).Int()
	if err != nil {
		return err
	}
	if capacity <= 0 {
		return ErrNoRequestCapacity
	}

	return repo.client.Decr(ctx, repo.endpoint+":"+KeyRequestCapacity).Err()
}

func (repo RedisCapacityRepository) updateDiskCapacity(ctx context.Context) error {
	stat := syscall.Statfs_t{}

	err := syscall.Statfs(repo.directoryPath, &stat)
	if err != nil {
		return err
	}

	return repo.client.Set(ctx, repo.endpoint+":"+KeyDiskCapacity, int64(stat.Bfree)*stat.Bsize, 0).Err()
}

func (repo RedisCapacityRepository) updateRequestCapacity(ctx context.Context) error {
	return repo.client.Set(ctx, repo.endpoint+":"+KeyRequestCapacity, repo.requestPerSecond, 0).Err()
}
