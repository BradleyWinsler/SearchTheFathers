package store

import (
	"context"

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
