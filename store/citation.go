package store

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

func (c *Client) GetCitations(ctx context.Context) ([]Citation, error) {
	filter := bson.M{}

	var cits []Citation

	curs, err := c.citationsColl.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = curs.Close(ctx)
	}()

	for curs.Next(ctx) {
		cit := Citation{}
		if err := curs.Decode(&cit); err != nil {
			return nil, err
		}
		cits = append(cits, cit)
	}
	if err := curs.Err(); err != nil {
		return nil, err
	}

	return cits, nil
}

func (c *Client) GetCitation(ctx context.Context, id string) (*Citation, error) {
	filter := bson.M{"_id": id}

	var cit *Citation

	res := c.citationsColl.FindOne(ctx, filter)
	if res.Err() != nil {
		return nil, res.Err()
	}

	if err := res.Decode(&cit); err != nil {
		return nil, err
	}

	return cit, nil
}

func (c *Client) InsertCitation(ctx context.Context, req *models.AddCitationRequest) (*Citation, error) {
	var tgs []store.Tag
	for _, t := range req.Tags {
		tgs = append(tgs, t)
	}

	ct := Citation{
		ID:                uuid.New(),
		Source:            req.Source,
		Father:            req.Father,
		Quote:             req.Quote,
		Tags:              tgs,
		Publisher:         req.Publisher,
		PublisherLocation: req.PublisherLocation,
		PublishDate:       req.PublishDate,
		Page:              req.Page,
		CreatedAt:         time.Now(),
		UpdatedAt:         0,
	}

	if _, err := c.citationsColl.InsertOne(ctx, &ct); err != nil {
		return nil, err
	}

	return &ct, nil
}

func (c *Client) DeleteCitation(ctx context.Context, id string) error {
	filter := bson.M{"_id": id}

	if _, err := c.citationsColl.DeleteOne(ctx, filter); err != nil {
		return err
	}

	return nil
}
