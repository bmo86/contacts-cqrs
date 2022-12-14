package repository

import (
	"contacts_cqrs/models"
	"context"
)

type RepositoryDB interface {
	Close()
	InsertCts(ctx context.Context, ct *models.Contact) error
	UpdateCts(ctx context.Context, id string, ct *models.Contact) error
	ListContact(ctx context.Context) ([]*models.Contact, error)
}

var repo RepositoryDB

func SetRepositoriDB(r RepositoryDB) {
	repo = r
}

func Close() {
	repo.Close()
}

func InsertCts(ctx context.Context, ct *models.Contact) error {
	return repo.InsertCts(ctx, ct)
}

func ListContact(ctx context.Context) ([]*models.Contact, error) {
	return repo.ListContact(ctx)
}

func UpdateCts(ctx context.Context, id string, ct *models.Contact) error {
	return repo.UpdateCts(ctx, id, ct)
}
