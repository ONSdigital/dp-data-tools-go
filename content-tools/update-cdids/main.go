package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ONSdigital/log.go/log"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	zebedeeURL  string
	mapperPath  string
	environment string
)

type CollectionDescription struct {
	Name            string      `json:"name"`
	Type            string      `json:"type"`
	PublishDate     interface{} `json:"publishDate"`
	Teams           []string    `json:"teams"`
	CollectionOwner string      `json:"collectionOwner"`
	ReleaseURI      interface{} `json:"releaseUri"`
	IsEncrypted     bool        `json:"isEncrypted"`
}

type CollectionDescriptionResponse struct {
	ID string `json:"id"`
}

func main() {
	setupFlags()

	ctx := context.Background()
	httpClient := http.DefaultClient

	if validateMandatoryParams(ctx) {
		return
	}

	collectionID, err := createCollection(ctx, httpClient, getCollectionName(), environment)
	if err != nil {
		fmt.Errorf("Error occurred while creating collection. Stopping the script. error: %s", err.Error())
		return
	}

	log.Event(ctx, "successfully updated all documents.", log.INFO)
}

func getCollectionName() string {
	return fmt.Sprintf("migrated-collection-%d", time.Now().Format(time.RFC3339))
}

func createCollection(ctx context.Context, client *http.Client, collectionName string, env string) (string, error) {

	//https://publishing.develop.onsdigital.co.uk/zebedee/collection
	collectionURL := fmt.Sprintf("%s/zebedee/collection", env)

	collectionDescription := &CollectionDescription{
		Name:            collectionName,
		Type:            "manual",
		PublishDate:     nil,
		Teams:           []string{"Test viewer team"},
		CollectionOwner: "ADMIN",
	}

	requestBodyString, err := json.Marshal(collectionDescription)
	if err != nil {
		errMessage := fmt.Errorf("failed to marshal data. Error: %v", err)
		log.Error(errMessage)
		return "", errMessage
	}

	creationResponse, err := client.Post(collectionURL, "application/json", bytes.NewBuffer(requestBodyString))
	if err != nil {
		errMessage := fmt.Errorf("failed to create collection via API. Error: %v", err)
		log.Error(errMessage)
		return "", errMessage
	}

	defer creationResponse.Body.Close()
	body, err := ioutil.ReadAll(creationResponse.Body)
	if err != nil {
		errMessage := fmt.Errorf("failed to read API response. Error: %v", err)
		log.Error(errMessage)
		return "", errMessage
	}
	response := &CollectionDescriptionResponse{}
	err = json.Unmarshal(body, response)
	if err != nil {
		errMessage := fmt.Errorf("failed to unmarshal API response. Error: %v", err)
		log.Error(errMessage)
		return "", errMessage
	}

	return response.ID, nil
}

func validateMandatoryParams(ctx context.Context) bool {
	if zebedeeURL == "" {
		log.Event(ctx, "missing zebedeeURL flag", log.ERROR)
		return true
	}

	if mapperPath == "" {
		log.Event(ctx, "missing mapper-path flag", log.ERROR)
		return true
	}

	if environment == "" {
		log.Event(ctx, "missing environment flag", log.ERROR)
		return true
	}
	return false
}

func setupFlags() {
	flag.StringVar(&zebedeeURL, "zebedee-url", zebedeeURL, "Zebedee API URL")
	flag.StringVar(&mapperPath, "mapper-path", mapperPath, "Path to the mapper")
	flag.StringVar(&environment, "environment-url", environment, "Environment URL")
	flag.Parse()
}
