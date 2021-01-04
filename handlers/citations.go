package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/BradleyWinsler/SearchTheFathers/models"
	"github.com/BradleyWinsler/SearchTheFathers/store"
)

func NewCitationHandlers(citationStore store.CitationStore) *CitationHandlers {
	return &CitationHandlers{
		citationStore: citationStore,
	}
}

type CitationHandlers struct {
	citationStore store.CitationStore
}

func (h *CitationHandlers) GetCitations(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	cits, err := h.citationStore.GetCitations(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(convertStoreToCitationModel(cits))
}

func (h *CitationHandlers) GetCitation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	cit, err := h.citationStore.GetCitation(context.Background(), params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c := convertStoreToCitationModel([]store.Citation{*cit})

	json.NewEncoder(w).Encoder(c[0])
}

func (h *CitationHandlers) AddCitation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var cr models.AddCitationRequest
	_ = json.NewDecoder(r.Body).Decode(&cr)

	if cr.Source == "" {
		http.Error(w, "Must include the source of this citation.", http.StatusBadRequest)
		return
	}
	if cr.Father == "" {
		http.Error(w, "Must include the father of this citation.", http.StatusBadRequest)
		return
	}
	if cr.Quote == "" {
		http.Error(w, "Must include the quote itself to add a citation.", http.StatusBadRequest)
		return
	}

	cit, err := h.citationStore.InsertCitation(context.Background(), &cr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c := convertStoreToCitationModel([]store.Citation{*cit})

	json.NewEncoder(w).Encode(c[0])
}

func (h *CitationHandlers) DeleteCitation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	if err := h.citationStore.DeleteCitation(context.Background(), params["id"]); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Successfully deleted citation."))
}

func convertStoreToCitationModel(cits []store.Citation) []models.Citation {
	var cs []models.Citation

	for _, c := range cits {
		cs = append(cs, models.Citation{
			ID:                c.ID,
			Source:            c.Source,
			Father:            c.Father,
			Quote:             c.Quote,
			Tags:              convertStoreToTagModel(c.Tags),
			Publisher:         c.Publisher,
			PublisherLocation: c.PublisherLocation,
			PublishDate:       c.PublishDate,
			Page:              c.Page,
			CreatedAt:         c.CreatedAt,
			UpdatedAt:         c.UpdatedAt,
		})
	}

	return cs
}

func convertStoreToTagModel(tags []store.Tag) []models.Tag {
	var ts []models.Tag

	for _, t := range tags {
		ts = append(ts, models.Tag{
			Slug: t.Slug,
		})
	}

	return ts
}
