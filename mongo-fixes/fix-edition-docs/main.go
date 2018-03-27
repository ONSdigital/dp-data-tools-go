package main

import (
	"errors"
	"flag"
	"strconv"
	"time"

	"github.com/ONSdigital/dp-dataset-api/models"
	"github.com/ONSdigital/go-ns/log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	mongoURL string
)

var states = []string{"edition-confirmed", "associated", "published"}

// ListCurrentEditions represents an object containing current edition resources
type ListCurrentEditions struct {
	Items []CurrentEdition `bson:"items"`
}

// CurrentEdition represents the current mongo edition document
type CurrentEdition struct {
	ID          string       `bson:"id,omitempty" json:"id,omitempty"`
	Edition     string       `bson:"edition" json:"edition"`
	Links       EditionLinks `bson:"links" json:"links"`
	State       string       `bson:"state" json:"state"`
	LastUpdated time.Time    `bson:"last_updated" json:"last_updated"`
}

// EditionLinks represents a list of edition links
type EditionLinks struct {
	Dataset       LinksObject `bson:"dataset" json:"dataset"`
	LatestVersion LinksObject `bson:"latest_version" json:"latest_version"`
	Self          LinksObject `bson:"self" json:"self"`
	Versions      LinksObject `bson:"versions" json:"versions"`
}

// LinksObject represents a generic links object
type LinksObject struct {
	HRef string `bson:"href,omitempty" json:"href,omitempty"`
	ID   string `bson:"id,omitempty" json:"id,omitempty"`
}

// NewEdition represents the next mongo edition document
type NewEdition struct {
	ID      string         `bson:"id" json:"id"`
	Current CurrentEdition `bson:"current,omitempty" json:"current,omitempty"`
	Next    CurrentEdition `bson:"next" json:"next"`
}

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

	// Get all editions
	editions, err := getEditions(session)
	if err != nil {
		log.ErrorC("failed to get all editions", err, nil)
		return
	}

	// loop over editions
	for _, edition := range editions.Items {
		logData := log.Data{"current_edition": edition}

		// Update edition doc to new structure
		newEdition := &NewEdition{
			ID:   edition.ID,
			Next: edition,
		}

		if edition.State == "published" {
			newEdition.Current = edition
		}

		// Get latest version for an edition (of state edition-confirmed, associated or published)
		state, version, err := getLatestVersion(session, edition.Links.Dataset.ID, edition.Edition, states)
		if err != nil {
			log.ErrorC("failed to get latest version", err, logData)
			return
		}

		// Update next sub document
		newEdition.Next.State = state
		newEdition.Next.Links.LatestVersion.ID = strconv.Itoa(version)
		newEdition.Next.Links.LatestVersion.HRef = "http://localhost:10400/datasets/" + edition.Links.Dataset.ID + "/editions/" + edition.Edition + "/versions/" + strconv.Itoa(version)

		if state != "published" {
			_, publishedVersion, err := getLatestVersion(session, edition.Links.Dataset.ID, edition.Edition, []string{"published"})
			if err != nil {
				if err != mgo.ErrNotFound {
					log.ErrorC("failed to get latest published version", err, logData)
					return
				}
			}

			newEdition.Current.Links.LatestVersion.ID = strconv.Itoa(publishedVersion)
			newEdition.Current.Links.LatestVersion.HRef = "http://localhost:10400/datasets/" + edition.Links.Dataset.ID + "/editions/" + edition.Edition + "/versions/" + strconv.Itoa(publishedVersion)
		} else {
			// If current data is wrong in mongo, fix it based on versions of the edition
			newEdition.Current.Links.LatestVersion.ID = strconv.Itoa(version)
			newEdition.Current.Links.LatestVersion.HRef = "http://localhost:10400/datasets/" + edition.Links.Dataset.ID + "/editions/" + edition.Edition + "/versions/" + strconv.Itoa(version)
		}

		logData["new_edition"] = newEdition

		// Remove current edition document
		if err = session.DB("datasets").C("editions").Remove(bson.M{"id": edition.ID}); err != nil {
			log.ErrorC("failed to delete edition before creating new edition", err, logData)
			return
		}

		// Create new Edition document
		if err = session.DB("datasets").C("editions").Insert(newEdition); err != nil {
			log.ErrorC("failed to insert new edition document, data lost in mongo but exists in this log", err, logData)
			return
		}

		log.Info("Successfully updated edition resource", logData)
	}
}

func getEditions(session *mgo.Session) (*ListCurrentEditions, error) {
	s := session.Copy()
	defer s.Close()

	iter := s.DB("datasets").C("editions").Find(bson.M{}).Iter()
	defer func() {
		err := iter.Close()
		if err != nil {
			log.ErrorC("error closing edition iterator", err, nil)
		}
	}()

	var results []CurrentEdition
	if err := iter.All(&results); err != nil {
		return nil, err
	}

	if len(results) < 1 {
		return nil, errors.New("no edition documents found")
	}

	return &ListCurrentEditions{Items: results}, nil
}

func getLatestVersion(session *mgo.Session, datasetID, edition string, listOfStates []string) (string, int, error) {
	s := session.Copy()
	defer s.Close()
	var version models.Version
	var latestVersion int

	selector := bson.M{
		"links.dataset.id": datasetID,
		"edition":          edition,
		"state": bson.M{
			"$in": listOfStates,
		},
	}

	// Results are sorted in reverse order to get latest version
	err := s.DB("datasets").C("instances").Find(selector).Sort("-version").One(&version)
	if err != nil {
		log.Info("We should never get here - this would mean there are no versions for a published edition assuming error is not found", log.Data{"dataset_id": datasetID, "edition": edition})
		return "", latestVersion, err
	}

	latestVersion = version.Version
	state := version.State

	return state, latestVersion, nil
}
