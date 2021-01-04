package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/BradleyWinsler/SearchTheFathers/store"
)

func NewTagHandlers(tagStore store.TagStore) *TagHandlers {
	return &TagHandlers{
		tagStore: tagStore,
	}
}

type TagHandlers struct {
	tagStore store.TagStore
}

func (h *TagHandlers) GetTags(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tags, err := h.tagStore.GetTags(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(convertStoreToTagModel(tags))
}

func (h *TagHandlers) AddTag(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	tag, err := h.tagStore.InsertTag(context.Background(), params["slug"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t := convertStoreToTagModel([]store.Tag{*tag})

	json.NewEncoder(w).Encode(t[0])
}

func (h *TagHandlers) DeleteTag(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	slug := params["slug"]

	if err := h.tagStore.DeleteTag(context.Background(), slug); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(fmt.Sprintf("Successfully deleted the %s tag.", slug)))
}
