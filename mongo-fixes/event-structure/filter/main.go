package main

import (
	"errors"
	"flag"
	"github.com/ONSdigital/go-ns/log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
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

	if mongoURL == "" {
		log.Error(errors.New("missing mongo-url flag"), nil)
		return
	}

	session, err := mgo.Dial(mongoURL)
	if err != nil {
		log.ErrorC("unable to create mongo session", err, nil)
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
			log.ErrorC("error closing edition iterator", err, nil)
		}
	}()

	errorCount := 0
	var filterIDs []string

	var filters []Filter
	if err := iter.All(&filters); err != nil {
		log.ErrorC("failed to get filters", err, nil)
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
		log.Info("failed to update all filter outputs.", log.Data{"number_of_unsuccessful_updates": errorCount, "filter_ids": filterIDs})
	}

	log.Info("successfully updated all documents.", nil)

}
