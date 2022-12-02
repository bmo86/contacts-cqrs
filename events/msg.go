package events

import "time"

type Message interface {
	Type() string
}

type CreatedContactMessage struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Lastname   string    `json:"lastname"`
	Image      string    `json:"image"`
	Email      string    `json:"email"`
	Phone      int       `json:"phone"`
	Status     bool      `json:"status"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdateData time.Time `json:"update_Data"`
}

func (m CreatedContactMessage) Type() string {
	return "Created_Contact"
}
