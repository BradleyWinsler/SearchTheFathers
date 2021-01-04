package store

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

func (c *Client) GetTags(ctx context.Context) ([]Tag, error) {
	filter := bson.M{}

	var tags []Tag

	curs, err := c.tagsColl.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = curs.Close(ctx)
	}()

	for curs.Next(ctx) {
		tag := Tag{}
		if err := curs.Decode(&tag); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	if err := curs.Err(); err != nil {
		return nil, err
	}

	return tags, nil
}

func (c *Client) InsertTag(ctx context.Context, slug string) (*Tag, error) {
	tg := Tag{
		Slug: slug,
	}

	if _, err := c.tagsColl.InsertOne(ctx, &tg); err != nil {
		return nil, err
	}

	return &tg, nil
}

func (c *Client) DeleteTag(ctx context.Context, slug string) error {
	filter := bson.M{"slug": slug}

	if _, err := c.tagsColl.DeleteOne(ctx, filter); err != nil {
		return err
	}

	return nil
}
