package main

import "time"

type CreatedContactMessage struct {
	ID         int       `json:"id"`
	Name       string    `json:"name" validate:"required"`
	Lastname   string    `json:"lastname" validate:"required"`
	Image      string    `json:"image" validate:"required"`
	Email      string    `json:"email" validate:"required"`
	Phone      int       `json:"phone" validate:"required"`
	Status     bool      `json:"status" validate:"required"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdateData time.Time `json:"update_Data"`
}

func newCreatedContactMessage(id int, name, lastname, image, email string, phone int, status bool, createdAt, updateDate time.Time) *CreatedContactMessage {
	return &CreatedContactMessage{
		ID:         id,
		Name:       name,
		Lastname:   lastname,
		Image:      image,
		Email:      email,
		Phone:      phone,
		Status:     status,
		CreatedAt:  createdAt,
		UpdateData: updateDate,
	}
}
