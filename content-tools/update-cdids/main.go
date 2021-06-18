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
	ZebedeeToken  = "ZEBEDEE_TOKEN"
	FlorenceToken = "X-Florence-Token"
)

type Config struct {
	Environment string
	Username    string
	Password    string
	FilePath    string
	Sheetname   string
	Limit       int64
}

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

type CdIDPair struct {
	cdID       string
	oldCdID    string
	ExcelIndex int
}

func (cp *CdIDPair) Print(ctx context.Context, prefix string) {
	log.Event(ctx, fmt.Sprintf("%s ExcelIndex: %d, cdid: %s, oldcdid: %s", prefix, cp.ExcelIndex, cp.cdID, cp.oldCdID), log.INFO)
}

func getNoticeMarkdown(newCdIDLocation string) string {
	return fmt.Sprintf("We have published a corrected format of this document [here][1].\n\n\n [1]: %s", newCdIDLocation)
}

func main() {
	config := parseConfig()

	ctx := context.Background()
	httpClient := http.DefaultClient

	if !config.isMandatoryParamsPresent(ctx) {
		return
	}

	cdIDPairs := readCdIDPairs(ctx, config)
	for _, pair := range cdIDPairs {
		pair.Print(ctx, "Parsed ")
	}
	log.Event(ctx, fmt.Sprintf("Read %d cdIDPairs", len(cdIDPairs)), log.INFO)
	ctx, err := authenticate(ctx, httpClient, config)
	if err != nil {
		errMessage := fmt.Errorf("error occurred while authenticating. Stopping the script. error: %s", err.Error())
		log.Error(errMessage)
		return
	}

	config.clearCreds()

	collectionID, err := createCollection(ctx, httpClient, getCollectionName(), config.Environment)
	if err != nil {
		errMessage := fmt.Errorf("error occurred while creating collection. Stopping the script. error: %s", err.Error())
		log.Error(errMessage)
		return
	}
	log.Event(ctx, fmt.Sprintf("created collection: %s", collectionID), log.INFO)

	for _, pair := range cdIDPairs {
		pair.Print(ctx, "Processing ")

		oldCdIDLocation, err := searchCdID(ctx, httpClient, config, pair.oldCdID)
		if err != nil {
			pair.Print(ctx, "stopping. Error occurred while fetching old cdID location")
			continue
		}

		newCDIDLocation, err := searchCdID(ctx, httpClient, config, pair.cdID)
		if err != nil {
			pair.Print(ctx, "stopping. Error occurred while fetching new cdID location")
			continue
		}

		existingCollectionID, err := checkIfCdIDExistsInAnotherCollection(ctx, httpClient, config, oldCdIDLocation)
		if err != nil {
			pair.Print(ctx, "stopping. Error occurred while verifying if cdID location exists in another collection")
			continue
		}

		if len(existingCollectionID) > 0 {
			pair.Print(ctx, fmt.Sprintf("stopping. Error occurred. CDID: %s already exists in another collection: %s", pair.oldCdID, existingCollectionID))
			continue
		}

		oldCdIDData, err := fetchDataForCDID(ctx, httpClient, config, oldCdIDLocation)
		if err != nil {
			pair.Print(ctx, "stopping. Error occurred while fetching cdID data.")
			continue
		}

		updatedCdIDData, err := addNoticeForNewCDID(ctx, oldCdIDData, pair, newCDIDLocation)
		if err != nil {
			pair.Print(ctx, "stopping. Error occurred while adding alert to cdidData")
			continue
		}

		err = addCDIDToCollection(ctx, httpClient, config, collectionID, oldCdIDLocation, updatedCdIDData)
		if err != nil {
			pair.Print(ctx, "stopping. Error occurred while adding CDID to collection")
			continue
		}

		err = approveCDID(ctx, httpClient, config, collectionID, oldCdIDLocation)
		if err != nil {
			pair.Print(ctx, "stopping. Error occurred while approving cdid in the collection")
			continue
		}

		pair.Print(ctx, "Completed Processing ")
	}

	log.Event(ctx, "successfully updated all documents.", log.INFO)
}

func checkIfCdIDExistsInAnotherCollection(ctx context.Context, client *http.Client, config *Config, cdIDLocation string) (string, error) {
	cdIDURI := strings.Replace(cdIDLocation, config.Environment, "", -1)
	collectionCheckURL := fmt.Sprintf("%s/checkcollectionsforuri?uri=%s/data.json", config.Environment, cdIDURI)

	req, err := http.NewRequestWithContext(ctx, "GET", collectionCheckURL, nil)
	if err != nil {
		errMessage := fmt.Errorf("failed to prepare search cdid in collection request. Error: %v", err)
		log.Error(errMessage)
		return "", errMessage
	}
	resp, err := client.Do(req)
	if err != nil {
		errMessage := fmt.Errorf("failed to search in API for cdid in collection . Error: %v", err)
		log.Error(errMessage)

		return "", errMessage
	}
	defer resp.Body.Close()

	if resp.StatusCode == 204 || resp.StatusCode > 400 {
		errMessage := fmt.Errorf("failed to read existing collection ID. Error In API: %v. Status Code: %s", err, resp.Status)
		log.Error(errMessage)

		return "", errMessage
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errMessage := fmt.Errorf("failed to read existing collection ID. Error: %v", err)
		log.Error(errMessage)
		return "", errMessage
	}

	return string(body), nil
}

// 	curl -X POST --header "X-Florwence-Token:$ZEBEDEE_TOKEN"
//	http://localhost:8082/review/{collection-id}?uri={path-to-content}/data.json&recursive=false
func approveCDID(ctx context.Context, client *http.Client, config *Config, collectionID string, cdIDLocation string) error {
	cdIDURI := strings.Replace(cdIDLocation, config.Environment, "", -1)

	pageReviewURL := fmt.Sprintf("%s/review/%s?uri=%s/data.json&recursive=false", config.Environment, collectionID, cdIDURI)
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
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
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

func addCDIDToCollection(ctx context.Context, client *http.Client, config *Config, collectionID string, cdIDLocation string, data string) error {
	cdIDURI := strings.Replace(cdIDLocation, config.Environment, "", -1)

	pageURL := fmt.Sprintf("%s/zebedee/content/%s?uri=%s/data.json&overwriteExisting=true", config.Environment, collectionID, cdIDURI)
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
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		errMessage := fmt.Errorf("failed to add collection to page. Error In API: %v. Status Code: %s", err, resp.Status)
		log.Error(errMessage)

		return errMessage
	}
	return nil
}

func fetchDataForCDID(ctx context.Context, client *http.Client, config *Config, location string) (string, error) {

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
	if resp.StatusCode != 200 {
		errMessage := fmt.Errorf("failed to read CDID. Error In API: %v. Status Code: %s", err, resp.Status)
		log.Error(errMessage)

		return "", errMessage
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errMessage := fmt.Errorf("failed to read CDID data response. Error: %v", err)
		log.Error(errMessage)
		return "", errMessage
	}

	return string(body), nil
}

func searchCdID(ctx context.Context, client *http.Client, config *Config, cdID string) (string, error) {
	searchURL := fmt.Sprintf("%s/search?q=%s", config.Environment, cdID)

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
	if resp.StatusCode != 200 {
		errMessage := fmt.Errorf("failed to search the given CDID. Error In API: %v. Status Code: %s", err, resp.Status)
		log.Error(errMessage)

		return "", errMessage
	}

	locationHeader := resp.Header.Get("Location")
	if len(locationHeader) == 0 {
		errMessage := fmt.Errorf("failed to search the given CDID")
		log.Error(errMessage)

		return "", errMessage
	}

	return strings.Split(locationHeader, "?")[0], nil
}

func authenticate(ctx context.Context, client *http.Client, config *Config) (context.Context, error) {
	loginURL := fmt.Sprintf("%s/login", config.Environment)

	credential := &Credential{
		Email:    config.Username,
		Password: config.Password,
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

	if loginResponse.StatusCode != 200 {
		errMessage := fmt.Errorf("failed to login. Error In API: %v. Status Code: %s", err, loginResponse.Status)
		log.Error(errMessage)

		return ctx, errMessage
	}

	body, err := ioutil.ReadAll(loginResponse.Body)
	if err != nil {
		errMessage := fmt.Errorf("failed to read API response. Error: %v", err)
		log.Error(errMessage)
		return ctx, errMessage
	}

	return context.WithValue(ctx, ZebedeeToken, string(body)), nil
}

func getCollectionName() string {
	return fmt.Sprintf("migrated-collection-%d", time.Now().UnixNano())
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

	if creationResponse.StatusCode != 200 {
		errMessage := fmt.Errorf("failed to create collection. Error In API: %v. Status Code: %s", err, creationResponse.Status)
		log.Error(errMessage)

		return "", errMessage
	}

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

func (config *Config) isMandatoryParamsPresent(ctx context.Context) bool {
	if config.Environment == "" {
		log.Event(ctx, "missing environment flag", log.ERROR)
		return false
	}

	if config.Username == "" {
		log.Event(ctx, "missing username flag", log.ERROR)
		return false
	}

	if config.Password == "" {
		log.Event(ctx, "missing password flag", log.ERROR)
		return false
	}

	if config.FilePath == "" {
		log.Event(ctx, "missing filepath flag", log.ERROR)
		return false
	}

	if config.Sheetname == "" {
		log.Event(ctx, "missing sheetname flag", log.ERROR)
		return false
	}

	if config.Limit == 0 {
		log.Event(ctx, "missing limit flag", log.ERROR)
		return false
	}
	return true
}

func (config *Config) clearCreds() {
	config.Username = ""
	config.Password = ""
}

func parseConfig() *Config {
	var environment string
	var username string
	var password string
	var filePath string
	var sheetname string
	var limit int64

	flag.StringVar(&environment, "environment-url", environment, "Environment URL")
	flag.StringVar(&password, "password", password, "password")
	flag.StringVar(&username, "username", username, "username")
	flag.StringVar(&filePath, "filepath", filePath, "filepath")
	flag.StringVar(&sheetname, "sheetname", sheetname, "sheetname to use")
	flag.Int64Var(&limit, "limit", limit, "limit of the cdids to process")
	flag.Parse()

	return &Config{
		Environment: environment,
		Username:    username,
		Password:    password,
		FilePath:    filePath,
		Sheetname:   sheetname,
		Limit:       limit,
	}
}

func readCdIDPairs(ctx context.Context, config *Config) []*CdIDPair {
	cdIDPairs := make([]*CdIDPair, 0)
	// Open given file.
	wb, err := xlsx.OpenFile(config.FilePath)
	if err != nil {
		panic(err)
	}
	// wb now contains a reference to the workbook
	// show all the sheets in the workbook

	log.Event(ctx, "Sheets in this file:", log.INFO)
	for i, sh := range wb.Sheets {
		log.Event(ctx, fmt.Sprintf("Index: %d Sheets name: %s in this file:", i, sh.Name), log.INFO)
	}

	sheet, ok := wb.Sheet[config.Sheetname]
	if !ok {
		panic(fmt.Errorf("sheet %s does not exist", config.Sheetname))
	}

	var rowCount int64
	for index, row := range sheet.Rows {
		cdIDPair := &CdIDPair{}
		if rowCount > config.Limit {
			break
		}
		if isValidRow(row) {
			cdIDPair.cdID = row.Cells[0].String()
			cdIDPair.oldCdID = row.Cells[3].String()
			cdIDPair.ExcelIndex = index
			cdIDPairs = append(cdIDPairs, cdIDPair)
		} else {
			log.Event(ctx, fmt.Sprintf("RowNumber: %d, Skipping cdid: %s, oldcdid: %s", index+1, row.Cells[0].String(), row.Cells[3].String()), log.INFO)
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

type Alert struct {
	Date     string `json:"date"`
	Markdown string `json:"markdown"`
	Type     string `json:"type"`
}

func addNoticeForNewCDID(ctx context.Context, dataJSON string, pair *CdIDPair, newCdIDLocation string) (string, error) {

	var result map[string]interface{}
	err := json.Unmarshal([]byte(dataJSON), &result)
	if err != nil {
		errMessage := fmt.Errorf("failed to unmarshal content data JSON response. Error: %v", err)
		log.Error(errMessage)
		return "", errMessage
	}
	existingAlerts := make([]Alert, 0)
	alerts, ok := result["alerts"].(map[string]interface{})
	if ok {
		for alert := range alerts {
			parsedAlert := Alert{}
			_ = json.Unmarshal([]byte(alert), &parsedAlert)
			existingAlerts = append(existingAlerts, parsedAlert)
		}
	}
	contentUpdateAlert := Alert{
		Date:     time.Now().Format(time.RFC3339),
		Markdown: getNoticeMarkdown(newCdIDLocation),
		Type:     "alert",
	}

	existingAlerts = append(existingAlerts, contentUpdateAlert)
	result["alerts"] = existingAlerts
	jsonString, err := json.Marshal(result)
	if err != nil {
		errMessage := fmt.Errorf("failed to marshal updated content. Error: %v", err)
		log.Error(errMessage)
		return "", errMessage
	}

	return string(jsonString), err
}
