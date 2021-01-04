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

func convertStoreToCitationModel(cits []store.Citation) []models.Citation {
	var cs []models.Citation

	for _, c := range cits {
		cs = append(cs, models.Citation{
			ID:                c.ID,
			Source:            c.Source,
			Father:            c.Father,
			Quote:             c.Quote,
			Tags:              c.Tags,
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
