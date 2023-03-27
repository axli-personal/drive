package adapters_test

import (
	"context"
	"github.com/axli-personal/drive/backend/storage/adapters"
	"testing"
	"time"
)

func TestRedisCapacityRepository(t *testing.T) {
	connectString := "redis://localhost:6379"

	repo, err := adapters.NewRedisCapacityRepository(connectString, "https://storage.example.com", directoryPath, 50)
	if err != nil {
		t.Fatal(err)
	}

	ticker := time.NewTicker(20 * time.Millisecond)
	deadline := time.After(10 * time.Second)

	for {
		select {
		case <-ticker.C:
			err = repo.DecreaseRequestCapacity(context.Background(), 1)
			if err != nil {
				t.Fatal(err)
			}
		case <-deadline:
			return
		}
	}
}
