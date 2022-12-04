package search

import (
	"bytes"
	"contacts_cqrs/models"
	"context"
	"encoding/json"
	"strconv"

	"github.com/elastic/go-elasticsearch/v7"
)

type ElasticRepo struct {
	client *elasticsearch.Client
}

func NewCNElastic(url string) (*ElasticRepo, error) {
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{url},
	})
	if err != nil {
		return nil, err
	}

	return &ElasticRepo{client: client}, nil
}

func (elas *ElasticRepo) Close() {
	//close jajajaj, hello
}

func (elas *ElasticRepo) SearchIndex(ctx context.Context, ct *models.Contact) error {
	body, _ := json.Marshal(ct)
	_, err := elas.client.Index(
		"contacts",
		bytes.NewReader(body),
		elas.client.Index.WithDocumentID(strconv.Itoa(ct.ID)),
		elas.client.Index.WithContext(ctx),
		elas.client.Index.WithRefresh("wait_for"),
	)
	return err
}
