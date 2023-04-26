package adapters_test

import (
	"bytes"
	"context"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/axli-personal/drive/backend/storage/adapters"
	"github.com/axli-personal/drive/backend/storage/domain"
	"io"
	"os"
	"path"
	"testing"
)

func TestDiskObjectRepository(t *testing.T) {
	testDir := path.Join(os.TempDir(), "drive-test")

	repo, err := adapters.NewDiskObjectRepository(testDir)
	if err != nil {
		t.Fatal(err)
	}

	content := bytes.Repeat([]byte("text content\n"), 100)
	data := bytes.NewReader(content)

	object, err := domain.NewObject(data)
	if err != nil {
		t.Fatal(err)
	}

	err = repo.SaveObject(context.Background(), object)
	if err != nil {
		t.Fatal(err)
	}

	ObjectInRepo, err := repo.GetObject(context.Background(), object.Hash())
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

func TestCloudObjectRepository(t *testing.T) {
	content := bytes.Repeat([]byte("text content\n"), 100)
	data := bytes.NewReader(content)

	object, err := domain.NewObject(data)
	if err != nil {
		t.Fatal(err)
	}

	endpoint := "oss-cn-hangzhou.aliyuncs.com"
	id := "LTAI5t7AVpkHzf81emXHdKFo"
	secret := "TZWv7aTpD9c5HhRUSmexqXRKb7VpqU"

	client, err := oss.New(endpoint, id, secret)
	if err != nil {
		t.Fatal(err)
	}

	bucket, err := client.Bucket("mintul")
	if err != nil {
		t.Fatal(err)
	}

	uploadSession, err := bucket.InitiateMultipartUpload(object.Hash())
	if err != nil {
		t.Fatal(err)
	}

	var parts []oss.UploadPart
	part, err := bucket.UploadPart(uploadSession, object, int64(object.TotalBytes()), 1)
	if err != nil {
		t.Fatal(err)
	}
	parts = append(parts, part)

	bucket.CompleteMultipartUpload(uploadSession, parts)
}
