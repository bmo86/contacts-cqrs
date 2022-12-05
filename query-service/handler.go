package main

import (
	"contacts_cqrs/events"
	"contacts_cqrs/means"
	"contacts_cqrs/models"
	"contacts_cqrs/repository"
	"contacts_cqrs/search"
	"context"
	"encoding/json"
	"log"
	"net/http"
)

func onCreatedContacts(m events.CreatedContactMessage) {
	ct := models.Contact{
		ID:         m.ID,
		Name:       m.Name,
		Lastname:   m.Lastname,
		Image:      m.Image,
		Email:      m.Email,
		Phone:      m.Phone,
		Status:     m.Status,
		CreatedAt:  m.CreatedAt,
		UpdateData: m.UpdateData,
	}

	if err := search.SearchIndex(context.Background(), &ct); err != nil {
		log.Printf("failed to index contact %v", err)
	}
}

func handlerList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error
	cts, err := repository.ListContact(ctx)
	if err != nil {
		means.MessageErr(http.StatusInternalServerError, err.Error(), w)
		return
	}

	w.Header().Set("Content-Type", "aplication/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cts)
}

func handlerSearch(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx := r.Context()
	query := r.URL.Query().Get("q")
	if len(query) == 0 {
		means.MessageErr(http.StatusBadRequest, "query is required", w)
		return
	}

	cts, err := search.SearchQuery(ctx, query)
	if err != nil {
		means.MessageErr(http.StatusInternalServerError, err.Error(), w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cts)
}
