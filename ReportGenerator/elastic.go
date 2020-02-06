package main

import (
	"log"
	"time"

	"context"

	"gopkg.in/olivere/elastic.v7"
)

// ES : Elasticsearch structure
type ES struct {
	client *elastic.Client
}

// Init connection to Elasticsearch
func (e *ES) Init(port string) (ok bool) {
	client, err := elastic.NewClient(elastic.SetURL(port))

	if err != nil {
		log.Println(err)
		return false
	}

	e.client = client
	return true
}

// GetDocumentsByQuery : Query only with carrier name and timeframe
func (e *ES) GetDocumentsByQuery(qs QueryBody, start time.Time, end time.Time) (hits *elastic.SearchResult, totalhits int64, ok bool) {

	// Create query
	query := elastic.NewBoolQuery()

	query.Filter(elastic.NewRangeQuery("timestamp").
		Format("strict_date_optional_time").
		From(start).
		To(end))

	if qs.CarrierName != "all" {
		query.Filter(elastic.NewMatchQuery("Carrier", qs.CarrierName))
	}

	if qs.Delayed != "all" {
		query.Filter(elastic.NewMatchQuery("FlightDelay", qs.Delayed))
	}

	if qs.Cancelled != "all" {
		query.Filter(elastic.NewMatchQuery("Cancelled", qs.Cancelled))
	}

	if qs.OriginCountry != "all" {
		query.Filter(elastic.NewMatchQuery("OriginCountry", qs.OriginCountry))
	}

	if qs.DestCountry != "all" {
		query.Filter(elastic.NewMatchQuery("DestCountry", qs.DestCountry))
	}

	// Perform the query
	res, err := e.client.Search().
		Index("flights").
		Type("doc").
		Query(query).
		Sort("timestamp", true).
		Size(10000).
		Do(context.Background())

	if err != nil {
		return nil, 0, false
	}

	return res, res.TotalHits(), true
}
