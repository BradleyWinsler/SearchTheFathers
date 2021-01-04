package store

type Citation struct {
	ID                string   `bson:"_id"`
	Source            string   `bson:"source"`
	Father            string   `bson:"father"`
	Quote             string   `bson:"quote"`
	Tags              []string `bson:"tags"`
	Publisher         string   `bson:"publisher"`
	PublisherLocation string   `bson:"publisher_location"`
	PublishDate       string   `bson:"publish_data"`
	Page              string   `bson:"page"`
}

type CitationsStore interface {
	GetCitations(ctx context.Context) ([]Citation, error)
}
