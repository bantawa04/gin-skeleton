package stripe

import (
	"context"
	"errors"
)

// ErrNotImplemented is returned until the Stripe SDK is wired into this placeholder.
var ErrNotImplemented = errors.New("stripe integration is not implemented")

// Config contains the minimum configuration expected by a Stripe implementation.
type Config struct {
	SecretKey string
}

// Client is a placeholder Stripe adapter.
type Client struct {
	config Config
}

// New creates a placeholder Stripe client.
func New(config Config) *Client {
	return &Client{config: config}
}

// CreatePaymentInput describes the data commonly needed to create a payment.
type CreatePaymentInput struct {
	Amount   int64
	Currency string
}

// CreatePayment is a placeholder for a Stripe payment implementation.
func (c *Client) CreatePayment(ctx context.Context, input CreatePaymentInput) error {
	_ = c
	_ = ctx
	_ = input
	return ErrNotImplemented
}
