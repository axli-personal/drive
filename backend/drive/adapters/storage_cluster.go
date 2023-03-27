package adapters

import (
	"context"
	"errors"
	"github.com/axli-personal/drive/backend/drive/domain"
	"github.com/axli-personal/drive/backend/drive/remote"
	"github.com/redis/go-redis/v9"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

const (
	MinDiskCapacity    = domain.MB
	KeyRequestCapacity = "RequestCapacity"
	KeyDiskCapacity    = "DiskCapacity"
)

var (
	ErrCannotChooseEndpoint = errors.New("cannot choose endpoint")
)

type Capacity struct {
	Disk    int
	Request int
}

type StorageScheduler struct {
	client      *redis.Client
	capacityMap map[string]Capacity
	lock        sync.Mutex
}

func NewStorageScheduler(connectionString string, endpoints []string) (remote.StorageCluster, error) {
	options, err := redis.ParseURL(connectionString)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(options)

	err = client.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}

	scheduler := &StorageScheduler{
		client:      client,
		capacityMap: make(map[string]Capacity),
	}

	for _, endpoint := range endpoints {
		scheduler.capacityMap[endpoint] = Capacity{}
	}

	go func() {
		ticker := time.NewTicker(time.Second)

		for range ticker.C {
			scheduler.updateCapacityMap(context.Background())
		}
	}()

	return scheduler, nil
}

func (scheduler *StorageScheduler) ChooseStorageEndPoint(ctx context.Context) (string, error) {
	scheduler.lock.Lock()
	defer scheduler.lock.Unlock()

	totalPriority := 0
	for _, capacity := range scheduler.capacityMap {
		if capacity.Request <= 0 {
			continue
		}
		if capacity.Disk < MinDiskCapacity {
			continue
		}

		totalPriority += capacity.Disk / MinDiskCapacity
	}

	if totalPriority == 0 {
		return "", ErrCannotChooseEndpoint
	}

	targetPriority := rand.Intn(totalPriority) + 1

	currentPriority := 0
	for endpoint, capacity := range scheduler.capacityMap {
		if capacity.Request <= 0 {
			continue
		}
		if capacity.Disk < MinDiskCapacity {
			continue
		}

		currentPriority += capacity.Disk / MinDiskCapacity

		if currentPriority >= targetPriority {
			return endpoint, nil
		}
	}

	return "", ErrCannotChooseEndpoint
}

func (scheduler *StorageScheduler) updateCapacityMap(ctx context.Context) error {
	scheduler.lock.Lock()
	defer scheduler.lock.Unlock()

	var lastErr error

	for endpoint := range scheduler.capacityMap {
		result, err := scheduler.client.Get(ctx, endpoint+":"+KeyRequestCapacity).Result()
		if err != nil {
			lastErr = err
			continue
		}

		requestCapacity, err := strconv.Atoi(result)
		if err != nil {
			lastErr = err
			continue
		}

		result, err = scheduler.client.Get(ctx, endpoint+":"+KeyDiskCapacity).Result()
		if err != nil {
			lastErr = err
			continue
		}

		diskCapacity, err := strconv.Atoi(result)
		if err != nil {
			lastErr = err
			continue
		}

		scheduler.capacityMap[endpoint] = Capacity{
			Request: requestCapacity,
			Disk:    diskCapacity,
		}
	}

	return lastErr
}
