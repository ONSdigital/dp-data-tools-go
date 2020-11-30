package main

import (
	"context"
	"flag"
	"time"

	"github.com/ONSdigital/log.go/log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Filter represents filter output resource
type Filter struct {
	FilterID string   `bson:"filter_id"`
	Events   []*Event `bson:"events,omitempty"     json:"events,omitempty"`
}

type Event struct {
	Type string    `bson:"type,omitempty" json:"type"`
	Time time.Time `bson:"time,omitempty" json:"time"`
}

func main() {

	var mongoURL string
	flag.StringVar(&mongoURL, "mongo-url", mongoURL, "mongoDB URL")
	flag.Parse()

	ctx := context.Background()

	if mongoURL == "" {
		log.Event(ctx, "missing mongo-url flag", log.ERROR)
		return
	}

	session, err := mgo.Dial(mongoURL)
	if err != nil {
		log.Event(ctx, "unable to create mongo session", log.ERROR, log.Error(err))
		return
	}
	defer session.Close()

	session.SetBatch(10000)
	session.SetPrefetch(0.25)

	collection := session.DB("filters").C("filterOutputs")

	iter := collection.Find(bson.M{}).Iter()
	defer func() {
		err := iter.Close()
		if err != nil {
			log.Event(ctx, "error closing edition iterator", log.ERROR, log.Error(err))
		}
	}()

	errorCount := 0
	var filterIDs []string

	var filters []Filter
	if err := iter.All(&filters); err != nil {
		log.Event(ctx, "failed to get filters", log.ERROR, log.Error(err))
		return
	}

	for _, filter := range filters {

		query := bson.M{
			"$set": bson.M{
				"events": []*Event{},
			},
		}

		if _, err := collection.Upsert(bson.M{"filter_id": filter.FilterID}, query); err != nil {
			filterIDs = append(filterIDs, filter.FilterID)
			errorCount++
		}
	}

	if errorCount > 0 {
		log.Event(ctx, "failed to update all filter outputs", log.ERROR, log.Data{"number_of_unsuccessful_updates": errorCount, "filter_ids": filterIDs})
	}

	log.Event(ctx, "successfully updated all documents", log.INFO)

}
