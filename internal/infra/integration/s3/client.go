package s3

import (
	"context"
	"errors"
	"io"
)

// ErrNotImplemented is returned until the AWS S3 SDK is wired into this placeholder.
var ErrNotImplemented = errors.New("s3 integration is not implemented")

// Config contains common S3 connection settings.
type Config struct {
	Region   string
	Bucket   string
	Endpoint string
}

// Client is a placeholder S3 storage adapter.
type Client struct {
	config Config
}

// New creates a placeholder S3 client.
func New(config Config) *Client {
	return &Client{config: config}
}

// Upload is a placeholder for uploading an object to S3-compatible storage.
func (c *Client) Upload(ctx context.Context, key string, body io.Reader) error {
	_ = c
	_ = ctx
	_ = key
	_ = body
	return ErrNotImplemented
}
