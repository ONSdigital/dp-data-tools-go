package main

import (
	"context"
	"errors"
	"flag"
	"strings"
	"time"

	"github.com/ONSdigital/dp-dataset-api/models"
	"github.com/ONSdigital/log.go/log"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

var (
	mongoURL string
)

// InstanceWithID represents instance with the additional _id from mongo
type InstanceWithID struct {
	ID                bson.ObjectId               `bson:"_id"`
	Alerts            *[]models.Alert             `bson:"alerts,omitempty"                      json:"alerts,omitempty"`
	CollectionID      string                      `bson:"collection_id,omitempty"               json:"collection_id,omitempty"`
	Dimensions        []models.Dimension          `bson:"dimensions,omitempty"                  json:"dimensions,omitempty"`
	Downloads         *models.DownloadList        `bson:"downloads,omitempty"                   json:"downloads,omitempty"`
	Edition           string                      `bson:"edition,omitempty"                     json:"edition,omitempty"`
	Events            *[]models.Event             `bson:"events,omitempty"                      json:"events,omitempty"`
	Headers           *[]string                   `bson:"headers,omitempty"                     json:"headers,omitempty"`
	ImportTasks       *models.InstanceImportTasks `bson:"import_tasks,omitempty"                json:"import_tasks"`
	InstanceID        string                      `bson:"id,omitempty"                          json:"id,omitempty"`
	LastUpdated       time.Time                   `bson:"last_updated,omitempty"                json:"last_updated,omitempty"`
	LatestChanges     *[]models.LatestChange      `bson:"latest_changes,omitempty"              json:"latest_changes,omitempty"`
	Links             *models.InstanceLinks       `bson:"links,omitempty"                       json:"links,omitempty"`
	ReleaseDate       string                      `bson:"release_date,omitempty"                json:"release_date,omitempty"`
	State             string                      `bson:"state,omitempty"                       json:"state,omitempty"`
	Temporal          *[]models.TemporalFrequency `bson:"temporal,omitempty"                    json:"temporal,omitempty"`
	TotalObservations *int                        `bson:"total_observations,omitempty"          json:"total_observations,omitempty"`
	UniqueTimestamp   bson.MongoTimestamp         `bson:"unique_timestamp"                      json:"-"`
	Version           int                         `bson:"version,omitempty"                     json:"version,omitempty"`
}

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
		log.Event(ctx, "unable to create mongo session", log.Error(err), log.ERROR)
		return
	}
	defer session.Close()

	session.SetBatch(10000)
	session.SetPrefetch(0.25)

	// Get all instances
	instances, err := getInstances(ctx, session)
	if err != nil {
		log.Event(ctx, "failed to get all instances", log.Error(err), log.ERROR)
		return
	}

	// Create a backup collection
	for _, instance := range instances {
		if err = addInstanceToBackup(session, instance); err != nil {
			log.Event(ctx, "failed to backup instances", log.Error(err), log.ERROR)
			return
		}
	}

	errorCount := 0

	// loop over instances
	for _, instance := range instances {

		// loop over dimensions
		for i, dimension := range instance.Dimensions {
			instance.Dimensions[i].HRef = strings.Replace(dimension.HRef, "/v1", "", 1)
		}

		// prepares updated_instance in bson.M and then updates existing instance document
		updatedInstance := bson.M{"$set": instance}

		// update the whole instance document
		if err := updateInstance(session, instance.InstanceID, updatedInstance); err != nil {
			log.Event(ctx, "failed to update dimension of instance", log.Error(err), log.Data{"current_instance": instance}, log.ERROR)
			errorCount++
		}

	}

	if errorCount > 0 {
		log.Event(ctx, "failed to update dimension of all instances", log.Data{"no_unsuccessful_updates": errorCount}, log.INFO)
	} else {
		log.Event(ctx, "successfully updated all instance documents", log.INFO)
	}

}

func getInstances(ctx context.Context, session *mgo.Session) (results []InstanceWithID, err error) {
	s := session.Copy()
	defer s.Close()

	iter := s.DB("datasets").C("instances").Find(bson.M{}).Iter()
	defer func() {
		err := iter.Close()
		if err != nil {
			log.Event(ctx, "error closing edition iterator", log.Error(err), log.ERROR)
		}
	}()

	if err := iter.All(&results); err != nil {
		return nil, err
	}

	if len(results) < 1 {
		return nil, errors.New("no instance documents found")
	}

	return results, nil
}

//createBackup updates an instance document
func addInstanceToBackup(session *mgo.Session, instance InstanceWithID) error {

	s := session.Copy()
	defer s.Close()

	err := s.DB("datasets").C("instances_backup").Insert(instance)
	return err
}

//UpdateInstance updates an instance document
func updateInstance(session *mgo.Session, instanceID string, updatedInstance bson.M) (err error) {
	s := session.Copy()
	defer s.Close()

	err = s.DB("datasets").C("instances").Update(bson.M{"id": instanceID}, updatedInstance)
	if err != nil {
		if err == mgo.ErrNotFound {
			return errors.New("instance not found")
		}
	}
	return err
}
