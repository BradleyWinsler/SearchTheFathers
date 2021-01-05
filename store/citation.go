package store

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/BradleyWinsler/SearchTheFathers/models"
)

func (c *Client) getCitations(ctx context.Context, filter bson.M) ([]Citation, error) {
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

func (c *Client) GetAllCitations(ctx context.Context) ([]Citation, error) {
	filter := bson.M{}

	cits, err := c.getCitations(ctx, filter)
	if err != nil {
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

func (c *Client) SearchCitations(ctx context.Context, req *models.SearchCitationsRequest) ([]Citation, error) {
	var cits []Citation
	hasTags := len(req.Tags) > 0

	if req.Source != "" {
		filter := bson.M{"source": req.Source}
		if cs, err := c.getCitations(ctx, filter); err != nil {
			return nil, err
		} else {
			return cs, nil
		}
	}

	if req.Father != "" {
		filter := bson.M{"father": req.Father}
		cs, err := c.getCitations(ctx, filter)
		if err != nil {
			return nil, err
		}
		if !hasTags {
			return cs, nil
		} else {
			return filterByTags(cs, req.Tags), nil
		}
	}

	if hasTags {
		if cs, err := c.GetAllCitations(ctx); err != nil {
			return nil, err
		} else {
			return filterByTags(cs, req.Tags), nil
		}
	}

	return cits, nil
}

func filterByTags(cits []Citation, tags []models.Tag) []Citation {
	var cs []Citation
	for _, c := range cits {
		for _, t := range c.Tags {
			for _, tg := range tags {
				if t.Slug == tg.Slug {
					cs = append(cs, c)
				}
			}
		}
	}

	return cs
}

func (c *Client) InsertCitation(ctx context.Context, req *models.AddCitationRequest) (*Citation, error) {
	var tgs []Tag
	for _, t := range req.Tags {
		tgs = append(tgs, Tag{
			Slug: t.Slug,
		})
	}

	ct := Citation{
		ID:                uuid.New().String(),
		Source:            req.Source,
		Father:            req.Father,
		Quote:             req.Quote,
		Tags:              tgs,
		Publisher:         req.Publisher,
		PublisherLocation: req.PublisherLocation,
		PublishDate:       req.PublishDate,
		Page:              req.Page,
		CreatedAt:         time.Now().UnixNano(),
		UpdatedAt:         0,
	}

	if _, err := c.citationsColl.InsertOne(ctx, &ct); err != nil {
		return nil, err
	}

	return &ct, nil
}

func (c *Client) UpdateCitation(ctx context.Context, id string, req *models.AddCitationRequest) (*Citation, error) {
	cit, err := c.GetCitation(ctx, id)
	if err != nil {
		return nil, err
	}

	ct := Citation{
		ID:        id,
		Tags:      cit.Tags,
		CreatedAt: cit.CreatedAt,
		UpdatedAt: time.Now().UnixNano(),
	}

	if req.Source != "" {
		ct.Source = req.Source
	}
	if req.Father != "" {
		ct.Father = req.Father
	}
	if req.Quote != "" {
		ct.Quote = req.Quote
	}
	if req.Publisher != "" {
		ct.Publisher = req.Publisher
	}
	if req.PublisherLocation != "" {
		ct.PublisherLocation = req.PublisherLocation
	}
	if req.PublishDate != "" {
		ct.PublishDate = req.PublishDate
	}
	if req.Page != "" {
		ct.Page = req.Page
	}

	if err := c.DeleteCitation(ctx, id); err != nil {
		return nil, err
	}

	if _, err := c.citationsColl.InsertOne(ctx, &ct); err != nil {
		return nil, err
	}

	return &ct, nil
}

func (c *Client) AddTagToCitation(ctx context.Context, id, slug string) error {
	filter := bson.M{"_id": id}
	update := bson.M{
		"$push": bson.M{
			"tags": bson.M{
				"slug": slug,
			},
		},
	}

	if _, err := c.citationsColl.UpdateOne(ctx, filter, update); err != nil {
		return err
	}

	return nil
}

func (c *Client) RemoveTagFromCitation(ctx context.Context, id, slug string) error {
	filter := bson.M{"_id": id}
	tag := []Tag{Tag{Slug: slug}}
	update := bson.M{
		"$pull": bson.M{
			"tags": bson.M{
				"$in": tag,
			},
		},
	}

	if _, err := c.citationsColl.UpdateOne(ctx, filter, update); err != nil {
		return err
	}

	return nil
}

func (c *Client) DeleteCitation(ctx context.Context, id string) error {
	filter := bson.M{"_id": id}

	if _, err := c.citationsColl.DeleteOne(ctx, filter); err != nil {
		return err
	}

	return nil
}
