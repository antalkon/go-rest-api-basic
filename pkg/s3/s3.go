package s3

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Config struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
}

type Client struct {
	Minio  *minio.Client
	Bucket string
}

// New инициализирует клиента MinIO
func New(cfg Config) (*Client, error) {
	minioClient, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("s3: failed to create client: %w", err)
	}

	// Проверим или создадим бакет
	ctx := context.Background()
	exists, err := minioClient.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		return nil, fmt.Errorf("s3: check bucket error: %w", err)
	}

	if !exists {
		err = minioClient.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("s3: create bucket error: %w", err)
		}
	}

	return &Client{
		Minio:  minioClient,
		Bucket: cfg.Bucket,
	}, nil
}
