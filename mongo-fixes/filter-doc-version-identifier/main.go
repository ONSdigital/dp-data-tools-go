package main

import (
	"context"
	"errors"
	"flag"
	"os"
	"time"

	"github.com/ONSdigital/dp-data-tools/mongo-fixes/filter-doc-version-identifier/data"
	"github.com/ONSdigital/dp-dataset-api/models"
	"github.com/ONSdigital/log.go/log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const publishedState = "published"

var (
	mongoURL string
)

func main() {
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

	count := 0

	count, err = updateFilters(ctx, session, "filters", count)
	if err != nil {
		log.Event(ctx, "failed updating filter blueprints", log.ERROR, log.Error(err))
		os.Exit(1)
	}
	log.Event(ctx, "successfuly updated filter blueprints", log.INFO)

	if _, err = updateFilters(ctx, session, "filterOutputs", count); err != nil {
		log.Event(ctx, "failed updating filter outputs", log.ERROR, log.Error(err))
		os.Exit(1)
	}

	log.Event(ctx, "successfuly updated filter outputs", log.INFO)
}

func updateFilters(ctx context.Context, session *mgo.Session, collection string, count int) (int, error) {
	logData := log.Data{"collection": collection}
	// Get all filters
	filters, err := getFilters(ctx, session, collection)
	if err != nil {
		log.Event(ctx, "failed to get all filter blueprints", log.ERROR, log.Error(err))
		return count, err
	}

	// loop over editions
	for _, filter := range filters.Items {
		count++
		if count%100 == 0 {

			log.Event(ctx, "getting existing count", log.INFO, log.Data{"count": count})
		}

		logData["filter"] = filter

		// Check dataset object does not already exists
		if filter.Dataset != nil && filter.Dataset.ID != "" {
			continue
		}

		// Get version, edition and dataset id for filter blueprint
		version, err := getVersion(ctx, session, filter.InstanceID)
		if err != nil {

			log.Event(ctx, "failed to get version doc", log.ERROR, log.Error(err), logData)
			return count, err
		}

		// Update filter blueprint document
		published := version.State == publishedState

		update := bson.M{
			"$set": bson.M{
				"dataset.id":      version.Links.Dataset.ID,
				"dataset.edition": version.Edition,
				"dataset.version": version.Version,
				"published":       published,
				"last_updated":    time.Now(),
			},
		}
		logData["update"] = update

		// Update filterBlueprint
		if err = session.DB("filters").C(collection).Update(bson.M{"filter_id": filter.FilterID}, update); err != nil {
			log.Event(ctx, "failed to update filter", log.ERROR, log.Error(err), logData)
			return count, err
		}
	}

	return count, nil
}

func getFilters(ctx context.Context, session *mgo.Session, collection string) (*data.ListFilters, error) {
	s := session.Copy()
	defer s.Close()

	iter := s.DB("filters").C(collection).Find(bson.M{}).Iter()
	defer func() {
		err := iter.Close()
		if err != nil {

			log.Event(ctx, "error closing filter iterator", log.ERROR, log.Error(err))
		}
	}()

	var results []data.Filter
	if err := iter.All(&results); err != nil {
		return nil, err
	}

	if len(results) < 1 {
		return nil, errors.New("no filter blueprints documents found")
	}

	return &data.ListFilters{Items: results}, nil
}

func getVersion(ctx context.Context, session *mgo.Session, instanceID string) (models.Version, error) {
	s := session.Copy()
	defer s.Close()
	var version models.Version

	// Results are sorted in reverse order to get latest version
	err := s.DB("datasets").C("instances").Find(bson.M{"id": instanceID}).One(&version)
	if err != nil {
		log.Event(ctx, "We should never get here - this would mean there are no versions for filter", log.ERROR, log.Error(err), log.Data{"instance_id": instanceID})
		return version, err
	}

	return version, nil
}
