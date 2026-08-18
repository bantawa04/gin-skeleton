package resend

import (
	"context"
	"errors"
)

// ErrNotImplemented is returned until the Resend SDK is wired into this placeholder.
var ErrNotImplemented = errors.New("resend integration is not implemented")

// Config contains the minimum configuration expected by a Resend implementation.
type Config struct {
	APIKey string
	From   string
}

// Client is a placeholder Resend email adapter.
type Client struct {
	config Config
}

// New creates a placeholder Resend client.
func New(config Config) *Client {
	return &Client{config: config}
}

// SendEmailInput describes a basic transactional email request.
type SendEmailInput struct {
	To      []string
	Subject string
	HTML    string
}

// SendEmail is a placeholder for sending an email through Resend.
func (c *Client) SendEmail(ctx context.Context, input SendEmailInput) error {
	_ = c
	_ = ctx
	_ = input
	return ErrNotImplemented
}
