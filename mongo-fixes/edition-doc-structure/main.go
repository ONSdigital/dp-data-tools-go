package main

import (
	"context"
	"errors"
	"flag"
	"strconv"
	"time"

	"github.com/ONSdigital/dp-dataset-api/models"
	"github.com/ONSdigital/log.go/log"
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

	// Get all editions
	editions, err := getEditions(ctx, session)
	if err != nil {
		log.Event(ctx, "failed to get all editions", log.ERROR, log.Error(err))
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
		state, version, err := getLatestVersion(ctx, session, edition.Links.Dataset.ID, edition.Edition, states)
		if err != nil {
			log.Event(ctx, "failed to get latest version", log.ERROR, log.Error(err), logData)
			return
		}

		// Update next sub document
		newEdition.Next.State = state
		newEdition.Next.Links.LatestVersion.ID = strconv.Itoa(version)
		newEdition.Next.Links.LatestVersion.HRef = "http://localhost:10400/datasets/" + edition.Links.Dataset.ID + "/editions/" + edition.Edition + "/versions/" + strconv.Itoa(version)

		if state != "published" {
			_, publishedVersion, err := getLatestVersion(ctx, session, edition.Links.Dataset.ID, edition.Edition, []string{"published"})
			if err != nil {
				if err != mgo.ErrNotFound {
					log.Event(ctx, "failed to get latest published version", log.ERROR, log.Error(err), logData)
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
			log.Event(ctx, "failed to delete edition before creating new edition", log.ERROR, log.Error(err), logData)
			return
		}

		// Create new Edition document
		if err = session.DB("datasets").C("editions").Insert(newEdition); err != nil {
			log.Event(ctx, "failed to insert new edition document, data lost in mongo but exists in this log", log.ERROR, log.Error(err), logData)
			return
		}

		log.Event(ctx, "successfully updated edition resource", log.INFO, logData)
	}
}

func getEditions(ctx context.Context, session *mgo.Session) (*ListCurrentEditions, error) {
	s := session.Copy()
	defer s.Close()

	iter := s.DB("datasets").C("editions").Find(bson.M{}).Iter()
	defer func() {
		err := iter.Close()
		if err != nil {
			log.Event(ctx, "error closing edition iterator", log.ERROR, log.Error(err))
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

func getLatestVersion(ctx context.Context, session *mgo.Session, datasetID, edition string, listOfStates []string) (string, int, error) {
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
		log.Event(ctx, "We should never get here - this would mean there are no versions for a published edition assuming error is not found",
			log.ERROR, log.Error(err), log.Data{"dataset_id": datasetID, "edition": edition})
		return "", latestVersion, err
	}

	latestVersion = version.Version
	state := version.State

	return state, latestVersion, nil
}
