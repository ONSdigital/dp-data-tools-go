package main

import (
	"errors"
	"flag"
	"os"
	"time"

	"github.com/ONSdigital/dp-data-tools/mongo-fixes/filter-doc-version-identifier/data"
	"github.com/ONSdigital/dp-dataset-api/models"
	"github.com/ONSdigital/go-ns/log"
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

	count := 0

	count, err = updateFilters(session, "filters", count)
	if err != nil {
		log.ErrorC("failed updating filter blueprints", err, nil)
		os.Exit(1)
	}
	log.Info("Successfuly updated filter blueprints", nil)

	if _, err = updateFilters(session, "filterOutputs", count); err != nil {
		log.ErrorC("failed updating filter outputs", err, nil)
		os.Exit(1)
	}
	log.Info("Successfuly updated filter outputs", nil)
}

func updateFilters(session *mgo.Session, collection string, count int) (int, error) {
	logData := log.Data{"collection": collection}
	// Get all filters
	filters, err := getFilters(session, collection)
	if err != nil {
		log.ErrorC("failed to get all filter blueprints", err, nil)
		return count, err
	}

	// loop over editions
	for _, filter := range filters.Items {
		count++
		if count%100 == 0 {
			log.Debug("", log.Data{"count": count})
		}

		logData["filter"] = filter

		// Check dataset object does not already exists
		if filter.Dataset != nil && filter.Dataset.ID != "" {
			continue
		}

		// Get version, edition and dataset id for filter blueprint
		version, err := getVersion(session, filter.InstanceID)
		if err != nil {
			log.ErrorC("failed to get version doc", err, logData)
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
			log.ErrorC("failed to update filter", err, logData)
			return count, err
		}
	}

	return count, nil
}

func getFilters(session *mgo.Session, collection string) (*data.ListFilters, error) {
	s := session.Copy()
	defer s.Close()

	iter := s.DB("filters").C(collection).Find(bson.M{}).Iter()
	defer func() {
		err := iter.Close()
		if err != nil {
			log.ErrorC("error closing filter iterator", err, nil)
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

func getVersion(session *mgo.Session, instanceID string) (models.Version, error) {
	s := session.Copy()
	defer s.Close()
	var version models.Version

	// Results are sorted in reverse order to get latest version
	err := s.DB("datasets").C("instances").Find(bson.M{"id": instanceID}).One(&version)
	if err != nil {
		log.Info("We should never get here - this would mean there are no versions for filter", log.Data{"instance_id": instanceID})
		return version, err
	}

	return version, nil
}
