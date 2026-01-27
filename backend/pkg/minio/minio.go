package minio

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/asliddinberdiev/eirsystem/config"
	"github.com/asliddinberdiev/eirsystem/pkg/logger"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Client struct {
	api *minio.Client
	log logger.Logger
}

func New(cfg *config.Minio, isDev bool, log logger.Logger) (*Client, error) {
	minioClient, err := minio.New(cfg.GetAddr(), &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.RootUser, cfg.RootPassword, ""),
		Secure: !isDev,
	})
	if err != nil {
		return nil, fmt.Errorf("minio connection error: %w", err)
	}

	client := &Client{api: minioClient, log: log}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for _, bucket := range cfg.Buckets {
		if err := client.ensureBucket(ctx, bucket.Name, bucket.Public); err != nil {
			return nil, err
		}
	}

	return client, nil
}

func (c *Client) ensureBucket(ctx context.Context, bucketName string, isPublic bool) error {
	exists, err := c.api.BucketExists(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("check bucket exists failed: %w", err)
	}
	if !exists {
		err = c.api.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("create bucket failed: %w", err)
		}
		c.log.Info("created bucket", logger.String("name", bucketName))
	}
	if isPublic {
		policy := fmt.Sprintf(
			`{
            "Version": "2012-10-17",
            "Statement": [
                {
                    "Effect": "Allow",
                    "Principal": {"AWS": ["*"]},
                    "Action": ["s3:GetObject"],
                    "Resource": ["arn:aws:s3:::%s/*"]
                }
            ]
        	}`,
			bucketName,
		)

		if err := c.api.SetBucketPolicy(ctx, bucketName, policy); err != nil {
			return fmt.Errorf("set bucket policy failed: %w", err)
		}
		c.log.Info("bucket is now PUBLIC", logger.String("name", bucketName))
	}

	return nil
}

func (c *Client) Upload(ctx context.Context, bucket, objectName string, reader io.Reader, size int64, contentType string) error {
	opts := minio.PutObjectOptions{
		ContentType: contentType,
	}
	if contentType == "" {
		opts.ContentType = "application/octet-stream"
	}

	_, err := c.api.PutObject(ctx, bucket, objectName, reader, size, opts)
	if err != nil {
		return fmt.Errorf("upload error: %w", err)
	}
	return nil
}

func (c *Client) GetLink(ctx context.Context, bucket, objectName string, expiry time.Duration, downloadName string) (string, error) {
	reqParams := make(url.Values)

	if downloadName != "" {
		reqParams.Set("response-content-disposition", fmt.Sprintf("attachment; filename=\"%s\"", downloadName))
	} else {
		reqParams.Set("response-content-disposition", "inline")
	}

	u, err := c.api.PresignedGetObject(ctx, bucket, objectName, expiry, reqParams)
	if err != nil {
		return "", fmt.Errorf("link generation error: %w", err)
	}
	return u.String(), nil
}

func (c *Client) Delete(ctx context.Context, bucket, objectName string) error {
	opts := minio.RemoveObjectOptions{GovernanceBypass: true}
	return c.api.RemoveObject(ctx, bucket, objectName, opts)
}
