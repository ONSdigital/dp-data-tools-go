package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ONSdigital/log.go/log"
	"github.com/tealeg/xlsx"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	ZebedeeToken = "ZEBEDEE_TOKEN"
)

var (
	zebedeeURL  string
	environment string
	username    string
	password    string
	filePath    string
	sheetname   string
	limit       int64
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

type Credential struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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

	cdIDPairs := readCdIDPairs(filePath, sheetname, limit)
	for _, pair := range cdIDPairs {
		fmt.Println("Processing for pairs")
		fmt.Printf("cdid: %s, oldCdid: %s", pair.cdID, pair.oldCdID)
		fmt.Println("---------------------")
	}

	ctx, err := authenticate(ctx, httpClient, environment, username, password)
	if err != nil {
		_ = fmt.Errorf("error occurred while authenticating. Stopping the script. error: %s", err.Error())
		return
	}

	collectionID, err := createCollection(ctx, httpClient, getCollectionName(), environment)
	if err != nil {
		fmt.Errorf("error occurred while creating collection. Stopping the script. error: %s", err.Error())
		return
	}
	fmt.Printf("created collection: %s", collectionID)

	log.Event(ctx, "successfully updated all documents.", log.INFO)
}

type CdIDPair struct {
	cdID    string
	oldCdID string
}

func authenticate(ctx context.Context, client *http.Client, environment, username string, password string) (context.Context, error) {
	loginURL := fmt.Sprintf("%s/login", environment)

	credential := &Credential{
		Email:    username,
		Password: password,
	}

	requestBodyString, err := json.Marshal(credential)
	if err != nil {
		errMessage := fmt.Errorf("failed to marshal data. Error: %v", err)
		log.Error(errMessage)
		return ctx, errMessage
	}

	req, err := http.NewRequestWithContext(ctx, "POST", loginURL, bytes.NewBuffer(requestBodyString))
	if err != nil {
		errMessage := fmt.Errorf("failed to prepare login request. Error: %v", err)
		log.Error(errMessage)
		return ctx, errMessage
	}

	req.Header.Set("Content-Type", "application/json")
	loginResponse, err := client.Do(req)
	if err != nil {
		errMessage := fmt.Errorf("failed to login via API. Error: %v", err)
		log.Error(errMessage)

		return ctx, errMessage
	}

	defer loginResponse.Body.Close()
	body, err := ioutil.ReadAll(loginResponse.Body)
	if err != nil {
		errMessage := fmt.Errorf("failed to read API response. Error: %v", err)
		log.Error(errMessage)
		return ctx, errMessage
	}

	return context.WithValue(ctx, ZebedeeToken, string(body)), nil
}

func getCollectionName() string {
	return fmt.Sprintf("migrated-collection-%d", time.Now().Format(time.RFC3339))
}

func createCollection(ctx context.Context, client *http.Client, collectionName string, env string) (string, error) {

	// https://publishing.develop.onsdigital.co.uk/zebedee/collection
	collectionURL := fmt.Sprintf("%s/zebedee/collection", env)

	collectionDescription := &CollectionDescription{
		Name:            collectionName,
		Type:            "manual",
		PublishDate:     nil,
		Teams:           []string{"Test viewer team"},
		CollectionOwner: "PUBLISHING_SUPPORT",
	}

	requestBodyString, err := json.Marshal(collectionDescription)
	if err != nil {
		errMessage := fmt.Errorf("failed to marshal data. Error: %v", err)
		log.Error(errMessage)
		return "", errMessage
	}
	req, err := http.NewRequestWithContext(ctx, "POST", collectionURL, bytes.NewBuffer(requestBodyString))
	if err != nil {
		errMessage := fmt.Errorf("failed to prepare collection creation request. Error: %v", err)
		log.Error(errMessage)
		return "", errMessage
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(ZebedeeToken, ctx.Value(ZebedeeToken).(string))

	creationResponse, err := client.Do(req)
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

	if environment == "" {
		log.Event(ctx, "missing environment flag", log.ERROR)
		return true
	}

	if username == "" {
		log.Event(ctx, "missing username flag", log.ERROR)
		return true
	}

	if password == "" {
		log.Event(ctx, "missing password flag", log.ERROR)
		return true
	}

	if filePath == "" {
		log.Event(ctx, "missing filepath flag", log.ERROR)
		return true
	}

	if sheetname == "" {
		log.Event(ctx, "missing sheetname flag", log.ERROR)
		return true
	}

	if limit == 0 {
		log.Event(ctx, "missing limit flag", log.ERROR)
		return true
	}
	return false
}

func setupFlags() {
	flag.StringVar(&zebedeeURL, "zebedee-url", zebedeeURL, "Zebedee API URL")
	flag.StringVar(&environment, "environment-url", environment, "Environment URL")
	flag.StringVar(&password, "password", password, "password")
	flag.StringVar(&username, "username", username, "username")
	flag.StringVar(&filePath, "filepath", filePath, "filepath")
	flag.StringVar(&sheetname, "sheetname", sheetname, "sheetname to use")
	flag.Int64Var(&limit, "limit", limit, "limit of the cdids to process")
	flag.Parse()
}

func readCdIDPairs(filePath string, sheetname string, limit int64) []*CdIDPair {
	cdIDPairs := make([]*CdIDPair, 0)
	// Open given file.
	wb, err := xlsx.OpenFile(filePath)
	if err != nil {
		panic(err)
	}
	// wb now contains a reference to the workbook
	// show all the sheets in the workbook
	fmt.Println("Sheets in this file:")
	for i, sh := range wb.Sheets {
		fmt.Println(i, sh.Name)
	}

	sheet, ok := wb.Sheet[sheetname]
	if !ok {
		fmt.Printf("Sheet %s does not exist", sheetname)
	}

	var rowCount int64
	for _, row := range sheet.Rows {
		cdIDPair := &CdIDPair{}
		if rowCount > limit {
			break
		}
		if isValidRow(row) {
			cdIDPair.cdID = row.Cells[0].String()
			cdIDPair.oldCdID = row.Cells[3].String()
			cdIDPairs = append(cdIDPairs, cdIDPair)
		}
		rowCount++
	}

	return cdIDPairs
}

func isValidRow(row *xlsx.Row) bool {
	isValid := row.Cells[0].Value != "" &&
		row.Cells[0].Value != "CDID" &&
		row.Cells[3].Value != "" &&
		row.Cells[3].Value != "no_cdid" &&
		row.Cells[3].Value != "old_cdid"
	return isValid
}
