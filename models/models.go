package models

type Citation struct {
	ID     string `json:"id"`
	Source string `json:"source"`
	Father string `json:"father"`
	Quote  string `json:"quote"`
	// Optionals
	Tags              []string `json:"tags"`
	Publisher         string   `json:"publisher"`
	PublisherLocation string   `json:"publisher_location"`
	PublishDate       string   `json:"publish_date"`
	Page              string   `json:"page"`
	CreatedAt         int64    `json:"created_at"`
	UpdatedAt         int64    `json:"updated_at"`
}

type Tag struct {
	Slug string `json:"slug"`
}

type AddCitationRequest struct {
	Source            string   `json:"source"`
	Father            string   `json:"father"`
	Quote             string   `json:"quote"`
	Tags              []string `json:"tags"`
	Publisher         string   `json:"publisher"`
	PublisherLocation string   `json:"publisher_location"`
	PublisherDate     string   `json:"publisher_date"`
	Page              string   `json:"page"`
}
