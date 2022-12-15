package main

import (
	"contacts_cqrs/events"
	"contacts_cqrs/means"
	"contacts_cqrs/models"
	"contacts_cqrs/repository"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
)

type createdContactRequest struct {
	Name     string `json:"name"`
	Lastname string `json:"lastname"`
	Image    string `json:"image"`
	Email    string `json:"email"`
	Phone    int    `json:"phone"`
	Status   bool   `json:"status"`
}

func createdContactHandler(w http.ResponseWriter, r *http.Request) {
	var req createdContactRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		means.MessageErr(http.StatusBadRequest, err.Error(), w)
		return
	}

	dat := time.Now().UTC()

	ct := models.Contact{
		ID:         1,
		Name:       req.Name,
		Lastname:   req.Lastname,
		Image:      req.Image,
		Email:      req.Email,
		Phone:      req.Phone,
		Status:     req.Status,
		CreatedAt:  dat,
		UpdateData: dat,
	}

	if err := repository.InsertCts(r.Context(), &ct); err != nil {
		means.MessageErr(http.StatusInternalServerError, err.Error(), w)
		return
	}

	if err := events.PublishCreatedContact(r.Context(), &ct); err != nil {
		log.Printf("failed to publish created contact event : %s", err)
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ct)
}

func updateContactHandler(w http.ResponseWriter, r *http.Request) {
	var req createdContactRequest

	id := r.URL.Query().Get("id")
	if len(id) == 0 {
		means.MessageErr(http.StatusBadRequest, "id not input", w)
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Fatal(err.Error())
	}

	dat := time.Now().UTC()
	ct := models.Contact{
		ID:         1,
		Name:       req.Name,
		Lastname:   req.Lastname,
		Image:      req.Image,
		Email:      req.Email,
		Phone:      req.Phone,
		Status:     req.Status,
		CreatedAt:  dat,
		UpdateData: dat,
	}

	valid := validator.New()

	if err := valid.Struct(ct); err != nil {
		means.MessageErr(http.StatusBadRequest, err.Error(), w)
		return
	}

	if err := repository.UpdateCts(r.Context(), id, &ct); err != nil {
		means.MessageErr(http.StatusInternalServerError, err.Error(), w)
		return
	}

	if err := events.PublishCreatedContact(r.Context(), &ct); err != nil {
		log.Printf("failed to publish created contact event : %s", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Contacto Actualizado")
}

func deleteContactHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if len(id) == 0 {
		means.MessageErr(http.StatusBadRequest, "id not found", w)
		return
	}

	if err := repository.DeleteCts(r.Context(), id); err != nil {
		means.MessageErr(http.StatusInternalServerError, err.Error(), w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Contacto Eliminado!")
}
