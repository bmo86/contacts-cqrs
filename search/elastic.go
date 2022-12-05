package search

import (
	"bytes"
	"contacts_cqrs/models"
	"context"
	"encoding/json"
	"errors"
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

func (elas *ElasticRepo) SearchQuery(ctx context.Context, query string) (result []models.Contact, err error) {
	var buf bytes.Buffer
	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":            query,
				"fields":           []string{"name", "lastname", "phone", "email"},
				"fuzziness":        3,
				"cutoff_frequency": 0.0001,
			},
		},
	}

	if err = json.NewEncoder(&buf).Encode(searchQuery); err != nil {
		return nil, err
	}

	res, err := elas.client.Search(
		elas.client.Search.WithContext(ctx),
		elas.client.Search.WithIndex("contacts"),
		elas.client.Search.WithBody(&buf),
		elas.client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := res.Body.Close(); err != nil {
			result = nil
		}
	}()

	if res.IsError() {
		return nil, errors.New("elastic error : " + res.String())
	}

	var eRes map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&eRes); err != nil {
		return nil, err
	}

	var cts []models.Contact
	for _, hit := range eRes["hits"].(map[string]interface{})["hits"].([]interface{}) {
		ct := models.Contact{}
		src := hit.(map[string]interface{})["_source"]
		marshal, err := json.Marshal(src)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(marshal, &ct); err == nil {
			cts = append(cts, ct)
		}
	}
	return cts, nil
}
