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
	"strings"
	"time"
)

const (
	ZebedeeToken = "ZEBEDEE_TOKEN"
	FlorenceToken = "X-Florence-Token"
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

	cdIDPairs := readCdIDPairs(ctx, filePath, sheetname, limit)
	for _, pair := range cdIDPairs {
		pair.Print(ctx, "Parsed ")
	}
	log.Event(ctx, fmt.Sprintf("Read %d cdIDPairs" , len(cdIDPairs)), log.INFO)
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
	log.Event(ctx, fmt.Sprintf("created collection: %s", collectionID), log.INFO)

	for _, pair := range cdIDPairs {
		pair.Print(ctx, "Processing ")

		oldCdIDLocation, err := searchOldCdID(ctx, httpClient, environment, pair.oldCdID)
		if err != nil {
			pair.Print(ctx, "stopping. Error occurred while fetching cdID location")
			continue
		}

		cdIDData, err := fetchDataForCDID(ctx, httpClient, environment, oldCdIDLocation)
		if err != nil {
			pair.Print(ctx, "stopping. Error occurred while fetching cdID data.")
			continue
		}

		err = addCDIDToCollection(ctx, httpClient, environment, collectionID, oldCdIDLocation, cdIDData)
		if err != nil {
			pair.Print(ctx, "stopping. Error occurred while adding CDID to collection")
			continue
		}

		err = approveCDID(ctx, httpClient, environment, collectionID, oldCdIDLocation)
		if err != nil {
			pair.Print(ctx, "stopping. Error occurred while approving cdid in the collection")
			continue
		}

		pair.Print(ctx, "Completed Processing ")
	}

	log.Event(ctx, "successfully updated all documents.", log.INFO)
}

// 	curl -X POST --header "X-Florence-Token:$ZEBEDEE_TOKEN"
//	http://localhost:8082/review/{collection-id}?uri={path-to-content}/data.json&recursive=false
func approveCDID(ctx context.Context, client *http.Client, environment string, collectionID string, cdIDLocation string) error {
	cdIDURI := strings.Replace(cdIDLocation, environment, "", -1)

	pageReviewURL := fmt.Sprintf("%s/review/%s?uri=%s/data.json&recursive=false", environment, collectionID, cdIDURI)
	req, err := http.NewRequestWithContext(ctx, "POST", pageReviewURL, nil)
	if err != nil {
		errMessage := fmt.Errorf("failed to approve CDID. Error: %v", err)
		log.Error(errMessage)
		return errMessage
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(FlorenceToken, ctx.Value(ZebedeeToken).(string))

	resp, err := client.Do(req)
	if err != nil {
		errMessage := fmt.Errorf("failed to approve CDID. Error: %v", err)
		log.Error(errMessage)

		return errMessage
	}
	if resp.StatusCode > 400 {
		errMessage := fmt.Errorf("failed to approve CDID. Error In API: %v. Status Code: %s", err, resp.Status)
		log.Error(errMessage)

		return errMessage
	}
	return nil
}

// POST http://localhost:8081/zebedee/content
//		/test-4c153d5de19c33be41b772b9d6c27dfc917bda71f2d8925e6b95e17da3a7d8cf
//		?uri=inflationandpriceindices/timeseries/mb55/mm22/data.json
//		&overwriteExisting=true

func addCDIDToCollection(ctx context.Context, client *http.Client, environment string, collectionID string, cdIDLocation string, data string) error {
	cdIDURI := strings.Replace(cdIDLocation, environment, "", -1)

	pageURL := fmt.Sprintf("%s/zebedee/content/%s?uri=%s/data.json&overwriteExisting=true", environment, collectionID, cdIDURI)
	req, err := http.NewRequestWithContext(ctx, "POST", pageURL, bytes.NewBuffer([]byte(data)))
	if err != nil {
		errMessage := fmt.Errorf("failed to page adding to collection request. Error: %v", err)
		log.Error(errMessage)
		return errMessage
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(ZebedeeToken, ctx.Value(ZebedeeToken).(string))

	resp, err := client.Do(req)
	if err != nil {
		errMessage := fmt.Errorf("failed to add collection to page. Error: %v", err)
		log.Error(errMessage)

		return errMessage
	}
	if resp.StatusCode > 400 {
		errMessage := fmt.Errorf("failed to add collection to page. Error In API: %v. Status Code: %s", err, resp.Status)
		log.Error(errMessage)

		return errMessage
	}
	return nil
}

func fetchDataForCDID(ctx context.Context, client *http.Client, environment string, location string) (string, error) {

	cdIDDataURL := fmt.Sprintf("%s/data", location)

	req, err := http.NewRequestWithContext(ctx, "GET", cdIDDataURL, nil)
	if err != nil {
		errMessage := fmt.Errorf("failed to prepare fetch CDID request. Error: %v", err)
		log.Error(errMessage)
		return "", errMessage
	}

	resp, err := client.Do(req)
	if err != nil {
		if err != nil {
			errMessage := fmt.Errorf("failed to fetch CDID data from API. Error: %v", err)
			log.Error(errMessage)

			return "", errMessage
		}
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errMessage := fmt.Errorf("failed to read CDID data response. Error: %v", err)
		log.Error(errMessage)
		return "", errMessage
	}

	return string(body), nil
}

func searchOldCdID(ctx context.Context, client *http.Client, environment string, cdID string) (string, error) {
	searchURL := fmt.Sprintf("%s/search?q=%s", environment, cdID)

	req, err := http.NewRequestWithContext(ctx, "GET", searchURL, nil)
	if err != nil {
		errMessage := fmt.Errorf("failed to prepare search request. Error: %v", err)
		log.Error(errMessage)
		return "", errMessage
	}
	resp, err := client.Transport.RoundTrip(req)
	if err != nil {
		if err != nil {
			errMessage := fmt.Errorf("failed to search API. Error: %v", err)
			log.Error(errMessage)

			return "", errMessage
		}
	}

	defer resp.Body.Close()
	locationHeader := resp.Header.Get("Location")
	if len(locationHeader) == 0 {
		errMessage := fmt.Errorf("failed to search the given CDID")
		log.Error(errMessage)

		return "", errMessage
	}

	return strings.Split(locationHeader, "?")[0], nil
}

type CdIDPair struct {
	cdID    string
	oldCdID string
}

func (cp *CdIDPair) Print(ctx context.Context, prefix string) {
	log.Event(ctx, fmt.Sprintf("%s cdid: %s, oldcdid: %s", prefix, cp.cdID, cp.oldCdID), log.INFO)
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

func readCdIDPairs(ctx context.Context, filePath string, sheetName string, limit int64) []*CdIDPair {
	cdIDPairs := make([]*CdIDPair, 0)
	// Open given file.
	wb, err := xlsx.OpenFile(filePath)
	if err != nil {
		panic(err)
	}
	// wb now contains a reference to the workbook
	// show all the sheets in the workbook

	log.Event(ctx, "Sheets in this file:", log.INFO)
	for i, sh := range wb.Sheets {
		log.Event(ctx, fmt.Sprintf("Index:  Sheets name:  in this file:", i, sh.Name), log.INFO)
	}

	sheet, ok := wb.Sheet[sheetName]
	if !ok {
		panic(fmt.Errorf("sheet %s does not exist", sheetName))
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
		} else {
			log.Event(ctx, fmt.Sprintf("Skipping cdid: %s, oldcdid: %s", row.Cells[0].String(), row.Cells[3].String()), log.INFO)
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
