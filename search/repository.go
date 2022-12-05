package search

import (
	"contacts_cqrs/models"
	"context"
)

type RepositorySearch interface {
	Close()
	SearchIndex(ctx context.Context, ct *models.Contact) error
	SearchQuery(ctx context.Context, query string) ([]models.Contact, error)
}

var elastic RepositorySearch

func SetRepoSearch(e RepositorySearch) {
	elastic = e
}

func Close() {
	elastic.Close()
}

func SearchIndex(ctx context.Context, ct *models.Contact) error {
	return elastic.SearchIndex(ctx, ct)
}

func SearchQuery(ctx context.Context, query string) ([]models.Contact, error) {
	return elastic.SearchQuery(ctx, query)
}
