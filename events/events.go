package events

import (
	"contacts_cqrs/models"
	"context"
)

type EventStore interface {
	Close()
	PublishCreatedContact(ctx context.Context, ct *models.Contact) error
	SubcribeCreatedContact(ctx context.Context) (<-chan CreatedContactMessage, error)
	OnCreatedContact(f func(CreatedContactMessage)) error
}

var events EventStore

func SetEvent(eventStore EventStore) {
	events = eventStore
}

func PublishCreatedContact(ctx context.Context, ct *models.Contact) error {
	return events.PublishCreatedContact(ctx, ct)
}

func SubcribeCreatedContact(ctx context.Context) (<-chan CreatedContactMessage, error) {
	return events.SubcribeCreatedContact(ctx)
}

func OnCreatedContact(f func(CreatedContactMessage)) error {
	return events.OnCreatedContact(f)
}

func Close() {
	events.Close()
}
	