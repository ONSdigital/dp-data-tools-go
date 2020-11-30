package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/ONSdigital/log.go/log"
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

	ctx := context.Background()

	if mongoURL == "" {
		log.Event(ctx, "missing mongo-url flag", log.ERROR)
		return
	}

	if downloadServiceURL == "" {
		log.Event(ctx, "missing download-service-url flag", log.ERROR)
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

	collection := session.DB("datasets").C("instances")

	iter := collection.Find(bson.M{}).Iter()
	defer func() {
		err := iter.Close()
		if err != nil {
			log.Event(ctx, "error closing edition iterator", log.ERROR, log.Error(err))
		}
	}()

	var instances []Instance
	if err := iter.All(&instances); err != nil {
		log.Event(ctx, "could not get instances from mongo", log.ERROR, log.Error(err))
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
			log.Event(ctx, "could not upsert document", log.ERROR, log.Error(err), log.Data{"instance_id": instance.ID, "query": query})
			return
		}

	}

	log.Event(ctx, "successfully updated all documents.", log.INFO)
}
