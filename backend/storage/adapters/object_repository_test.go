package adapters_test

import (
	"bytes"
	"context"
	"github.com/axli-personal/drive/backend/storage/adapters"
	"github.com/axli-personal/drive/backend/storage/domain"
	"github.com/google/uuid"
	"io"
	"testing"
)

const directoryPath = "./data"

func TestDiskObjectRepository(t *testing.T) {
	repo, err := adapters.NewDiskObjectRepository(directoryPath)
	if err != nil {
		t.Fatal(err)
	}

	id := uuid.New().String()
	content := bytes.Repeat([]byte("text content\n"), 100)
	data := bytes.NewReader(content)

	object, err := domain.NewObject(id, data)
	if err != nil {
		t.Fatal(err)
	}

	err = repo.SaveObject(context.Background(), object)
	if err != nil {
		t.Fatal(err)
	}

	ObjectInRepo, err := repo.GetObject(context.Background(), id)
	if err != nil {
		t.Fatal(err)
	}

	contentInRepo, err := io.ReadAll(ObjectInRepo)
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Compare(content, contentInRepo) != 0 {
		t.Fatal("content not equal")
	}
}
