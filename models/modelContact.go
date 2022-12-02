package models

import "time"

type Contact struct {
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
