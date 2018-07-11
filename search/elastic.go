package search

import (
	"context"
	"encoding/json"
	"log"

	"github.com/olivere/elastic"
	"github.com/bialas1993/go-notifier/schema"
)

const ELASTIC_INDEX = "notifications"

type ElasticRepository struct {
	client *elastic.Client
}

func NewElastic(url string) (*ElasticRepository, error) {
	client, err := elastic.NewClient(
		elastic.SetURL(url),
		elastic.SetSniff(false),
	)
	if err != nil {
		return nil, err
	}
	return &ElasticRepository{client}, nil
}

func (r *ElasticRepository) Close() {
}

func (r *ElasticRepository) Insert(ctx context.Context, notify schema.Notify) error {
	_, err := r.client.Index().
		Index(ELASTIC_INDEX).
		Type("notify").
		Id(notify.ID).
		BodyJson(notify).
		Refresh("wait_for").
		Do(ctx)
	return err
}

func (r *ElasticRepository) Search(ctx context.Context, query string, skip uint64, take uint64) ([]schema.Notify, error) {
	result, err := r.client.Search().
		Index(ELASTIC_INDEX).
		Query(
			elastic.NewMultiMatchQuery(query, "body", "title").
				Fuzziness("3").
				PrefixLength(1).
				CutoffFrequency(0.0001),
		).
		From(int(skip)).
		Size(int(take)).
		Do(ctx)
	if err != nil {
		return nil, err
	}
	notifications := []schema.Notify{}
	for _, hit := range result.Hits.Hits {
		var notify schema.Notify
		if err = json.Unmarshal(*hit.Source, &notify); err != nil {
			log.Println(err)
		}
		notifications = append(notifications, notify)
	}
	return notifications, nil
}
