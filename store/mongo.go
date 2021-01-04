package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/cenkalti/backoff"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connection settings
const (
	maxRetries              = 5
	connectTimeoutInSeconds = 10
	pingTimeoutInSeconds    = 2
)

// Database & collections
const (
	database            = "fathers"
	citationsCollection = "citations"
)

// Client holds the dependencies needed to implement the store.
type Client struct {
	client        *mongo.Client
	citationsColl *mongo.Collection
}

// Setup and teardown
func NewClient(ctx context.Context, uri string) (*Client, error) {
	var client *mongo.Client
	if err := backoff.Retry(connect(uri, &client), backoff.WithMaxRetries(backoff.NewExponentialBackoff(), maxRetries)); err != nil {
		return nil, fmt.Errorf("error connecting to mongo database: %w", err)
	}
	if err := backoff.Retry(ping(client), backoff.WithMaxRetries(backoff.NewExponentialBackoff(), maxRetries)); err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	db := client.Database(database)

	var errs []error

	errs = append(errs, db.CreateCollection(ctx, citationsCollection))
	for _, err := range errs {
		if err != nil && !errors.As(err, &mongo.CommandError{}) {
			return nil, err
		}
	}

	citationsColl := db.Collection(citationsCollection)

	return &Client{
		client:        client,
		citationsColl: citationsColl,
	}, nil
}

// Close closes the mongo client
func Close(c *Client) error {
	return c.client.Disconnect(context.Background())
}

func connect(uri string, client **mongo.Client) backoff.Operation {
	return func() error {
		ctx, cancel := context.WithTimeout(context.Background(), connectTimeoutInSeconds*time.Second)
		defer cancel()

		c, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
		if err != nil {
			return err
		}

		// Assign to the client value
		*client = c

		return nil
	}
}

func ping(client *mongo.Client) backoff.Operation {
	return func() error {
		ctx, cancel := context.WithTimeout(context.Background(), pingTimeoutInSeconds*time.Second)
		defer cancel()

		if err := client.Ping(ctx, readpref.Primary()); err != nil {
			return err
		}
		return nil
	}
}
