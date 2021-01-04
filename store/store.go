package store

import (
	"context"
)

type Citation struct {
	ID                string   `bson:"_id"`
	Source            string   `bson:"source"`
	Father            string   `bson:"father"`
	Quote             string   `bson:"quote"`
	Tags              []string `bson:"tags"`
	Publisher         string   `bson:"publisher"`
	PublisherLocation string   `bson:"publisher_location"`
	PublishDate       string   `bson:"publish_date"`
	Page              string   `bson:"page"`
	CreatedAt         int64    `bson:"created_at"`
	UpdatedAt         int64    `bson:"updated_at"`
}

type CitationsStore interface {
	GetCitations(ctx context.Context) ([]Citation, error)
}
