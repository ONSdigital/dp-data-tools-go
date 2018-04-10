package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/ONSdigital/go-ns/log"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Filter represents filter output resource
type Filter struct {
	Downloads map[string]Download `bson:"downloads"`
	Dataset   *Dataset            `bson:"dataset"`
	FilterID  string              `bson:"filter_id"`
}

// Download represents download of filter
type Download struct {
	URL string `bson:"url"`
}

// Dataset represents a list of identifiers to find the unique version
// resource in which the filter output is filtering against
type Dataset struct {
	ID      string `bson:"id"`
	Edition string `bson:"edition"`
	Version int    `bson:"version"`
}

var (
	mongoURL           string
	downloadServiceURL string
)

func main() {
	flag.StringVar(&mongoURL, "mongo-url", mongoURL, "mongoDB URL")
	flag.StringVar(&downloadServiceURL, "download-service-url", downloadServiceURL, "download-service url")
	flag.Parse()

	if mongoURL == "" {
		log.Error(errors.New("missing mongo-url flag"), nil)
		return
	}

	if downloadServiceURL == "" {
		log.Error(errors.New("missing download-service-url flag"), nil)
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
		log.ErrorC("could not get instances from mongo", err, nil)
		return
	}

	for _, filter := range filters {

		href := fmt.Sprintf("%s/downloads/filter-outputs/%s", downloadServiceURL, filter.FilterID)

		query := bson.M{
			"$set": bson.M{
				"downloads.csv.href":   href + ".csv",
				"downloads.csv.public": filter.Downloads["csv"].URL,
				"downloads.xls.href":   href + ".xlsx",
				"downloads.xls.public": filter.Downloads["xls"].URL,
			},
		}

		if _, err := collection.Upsert(bson.M{"filter_id": filter.FilterID}, query); err != nil {
			filterIDs = append(filterIDs, filter.FilterID)
			errorCount++
		}
	}

	if errorCount > 0 {
		log.Info("unsuccessfully updated all documents.", log.Data{"number_of_unsuccessful_updates": errorCount, "filter_ids": filterIDs})
	}

	log.Info("successfully updated all documents.", nil)

}
