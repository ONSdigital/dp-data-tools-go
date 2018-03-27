package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/ONSdigital/go-ns/log"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Instance represents dataset instance
type Instance struct {
	Downloads map[string]Download `bson:"downloads"`
	ID        string              `bson:"id"`
	Version   int                 `bson:"version"`
	Edition   string              `bson:"edition"`
	Links     Links               `bson:"links"`
}

// Download represents download on instance
type Download struct {
	URL string `bson:"url"`
}

// Links represents links
type Links struct {
	Dataset Dataset `bson:"dataset"`
}

// Dataset represents dataset link
type Dataset struct {
	ID string `bson:"id"`
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

	collection := session.DB("datasets").C("instances")

	iter := collection.Find(bson.M{}).Iter()
	defer func() {
		err := iter.Close()
		if err != nil {
			log.ErrorC("error closing edition iterator", err, nil)
		}
	}()

	var instances []Instance
	if err := iter.All(&instances); err != nil {
		log.ErrorC("could not get instances from mongo", err, nil)
		return
	}

	for _, instance := range instances {

		href := fmt.Sprintf("%s/downloads/datasets/%s/editions/%s/versions/%d", downloadServiceURL, instance.Links.Dataset.ID, instance.Edition, instance.Version)

		query := bson.M{
			"$set": bson.M{
				"downloads.csv.href":   href + ".csv",
				"downloads.csv.public": instance.Downloads["csv"].URL,
				"downloads.xls.href":   href + ".xlsx",
				"downloads.xls.public": instance.Downloads["xls"].URL,
			},
		}

		if _, err := collection.Upsert(bson.M{"id": instance.ID}, query); err != nil {
			log.ErrorC("could not upsert document", err, log.Data{"instance_id": instance.ID, "query": query})
			return
		}

	}

	log.Info("successfully updated all documents.", nil)

}
