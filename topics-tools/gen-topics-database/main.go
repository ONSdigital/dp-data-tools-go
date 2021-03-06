package main

// NOTE: if vscode graphviz view can't cope with the size of the generated t.gv file
//       then use the graphviz command line command to convert to .svg:
//
//       dot -Tsvg t.gv -o t.gv.svg
//
//    OR: (if flag: GraphAllContent is true)
//       dot -Tsvg t-big.gv -o t-big.gv.svg
//
//    OR: (if flags: GraphAllContent and ColourContent are both true)
//       dot -Tsvg t-big-colour.gv -o t-big-colour.gv.svg

// NOTE: Any changes to the ONS website might stop this App working and it will need adjusting ..
//       That is, the struct(s) may need additional fields.

// NOTE: to grab all output info from running this use:
//       go run main.go topic.go content.go >t.txt

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config represents config for ons-scrape
type Config struct {
	ScrapeContent bool `envconfig:"SCRAPE_CONTENT"`
	// WARNING: if 'GraphAllContent' is set true, see NOTE at start of file about not being able to view graphviz in vscode
	GraphAllContent    bool `envconfig:"GRAPH_ALL_CONTENT"`
	ColourContent      bool `envconfig:"COLOUR_CONTENT"`        // this acts as an option on flag 'GraphAllContent'
	FullDepth          bool `envconfig:"FULL_DEPTH"`            // search recursively through the whole site from content down (this covers a lot of the site but may miss something)
	OnlyFirstFullDepth bool `envconfig:"ONLY_FIRST_FULL_DEPTH"` // enable this to minimise the amount of full depth recursions for code development / testing
	SkipVersions       bool `envconfig:"SKIP_VERSIONS"`         // when doing FULL_DEPTH, this skips processing of version files (to save time when developing this code)
	PlayNice           bool `envconfig:"PLAY_NICE"`             // add a little delay before reading each page
}

var cfg *Config

// Sensible combinations of flags are:
//
// 1. ScrapeContent: true  [the rest false]
//    This generates scripts: content-init.js and topic-init.js to set up mongo topics database
//    for the initial dp-topics-api project. (in directory: mongo-init-scripts)
//    This also generates:
//    t.gv : a graphviz document that can be viewed in vscode with the extension:
//           Graphviz (dot) language support for Visual Studio Code, by: João Pinto
//           (or see notes at top of this file)
//           (in directory: graphviz-files)
//    broken_links.txt : a report on links in ONS site that are not working
//           (in directory: observations)
//    If a page has problems being decoded into a go struct two files may get written
//    into directory 'temp' depending on the nature of the problem.
//    These two .json files can be formatted in vscode, saved and then visually compared
//    in an application like 'meld' to determine how to adjust go struct to decode page /data.
//    There are a number of other files generated for diagnostics and future development,
//    see the main() for functions called to see what these are (they are in a state of flux
//    and may change).
//
// 2. ScrapeContent: true, GraphAllContent: true, ColourContent: true
//    This causes the graphviz output file to become t-big-colour.gv and to view this it will
//    need to be converted into a .svg file with command:
//           dot -Tsvg t-big-colour.gv -o t-big-colour.gv.svg
//
// 3. ScrapeContent: true, FullDepth: true
//    This will scrape through the whole ONS site and generate more mongo init scripts named by
//    the 'type' of the page found.
//    This can take 30 minutes, to 5 hours to run, depending on which part of ONS site is scanned.
//    Also, the broken_links.txt file may be a lot bigger.
//
// 4. Other flags: OnlyFirstFullDepth, SkipVersions
//    When 'FullDepth' is true, setting these flags 'true' will reduce the amount of the ONS site
//    that is scanned. This is useful when developing the code and the definitions of the struct's.
//

// InitConfig returns the default config with any modifications through environment
// variables
func InitConfig() error {
	cfg = &Config{
		ScrapeContent:      true,
		GraphAllContent:    false,
		ColourContent:      false,
		FullDepth:          false,
		OnlyFirstFullDepth: false,
		SkipVersions:       false,
		PlayNice:           false,
	}

	return envconfig.Process("", cfg)
}

// NOTE: the order that fields appear within structs is critical for this App to work !

// DataResponse is the whole data page that gets filled from Unmarshal'ing pages like:
// https://www.ons.gov.uk/businessindustryandtrade/data
// and this is a top level 'node' like page
type DataResponse struct {
	Sections                  *[]SubLink                        `bson:"sections,omitempty"                   json:"sections,omitempty"`
	Items                     *[]ItemLinks                      `bson:"items,omitempty"                      json:"items,omitempty"`
	Datasets                  *[]DatasetLinks                   `bson:"datasets,omitempty"                   json:"datasets,omitempty"`
	HighlightedLinks          *[]HighlightLinks                 `bson:"highlightedLinks,omitempty"           json:"highlightedLinks,omitempty"`
	StatsBulletins            *[]StatsBulletLinks               `bson:"statsBulletins,omitempty"             json:"statsBulletins,omitempty"`
	RelatedArticles           *[]RelatedArticleLinks            `bson:"relatedArticles,omitempty"            json:"relatedArticles,omitempty"`
	RelatedMethodology        *[]RelatedMethodLinks             `bson:"relatedMethodology,omitempty"         json:"relatedMethodology,omitempty"`
	RelatedMethodologyArticle *[]RelatedMethodologyArticleLinks `bson:"relatedMethodologyArticle,omitempty"  json:"relatedMethodologyArticle,omitempty"`
	HighlightedContent        *[]HighlightedContentLinks        `bson:"highlightedContent,omitempty"         json:"highlightedContent,omitempty"`
	Index                     *int                              `bson:"index,omitempty"                      json:"index,omitempty"`
	Type                      *string                           `bson:"type,omitempty"                       json:"type,omitempty"`
	URI                       *string                           `bson:"uri,omitempty"                        json:"uri,omitempty"`
	Description               *Descript                         `bson:"description,omitempty"                json:"description,omitempty"`
}

// SubLink is the sub page links
type SubLink struct {
	URI   *string `bson:"uri,omitempty"    json:"uri,omitempty"`
	Index *int    `bson:"index,omitempty"  json:"index,omitempty"`
}

// HighlightLinks are highlights
type HighlightLinks struct {
	URI      *string `bson:"uri,omitempty"  json:"uri,omitempty"`
	valid    bool    // not read in or saved, but used by this App
	linkType string  // not read in or saved, but used by this App
}

// Descript is the page description
type Descript struct {
	Title           *string   `bson:"title,omitempty"            json:"title,omitempty"`
	Summary         *string   `bson:"summary,omitempty"          json:"summary,omitempty"`
	Keywords        *[]string `bson:"keywords,omitempty"         json:"keywords,omitempty"`
	MetaDescription *string   `bson:"metaDescription,omitempty"  json:"metaDescription,omitempty"`
	Unit            *string   `bson:"unit,omitempty"             json:"unit,omitempty"`
	PreUnit         *string   `bson:"preUnit,omitempty"          json:"preUnit,omitempty"`
	Source          *string   `bson:"source,omitempty"           json:"source,omitempty"`
}

// ItemLinks are highlights
type ItemLinks struct {
	URI      *string `bson:"uri,omitempty"  json:"uri,omitempty"`
	valid    bool    // not read in or saved, but used by this App
	linkType string  // not read in or saved, but used by this App
}

// DatasetLinks are highlights
type DatasetLinks struct {
	URI      *string `bson:"uri,omitempty"  json:"uri,omitempty"`
	valid    bool    // not read in or saved, but used by this App
	linkType string  // not read in or saved, but used by this App
}

// StatsBulletLinks are stats bulletins
type StatsBulletLinks struct {
	URI      *string `bson:"uri,omitempty"  json:"uri,omitempty"`
	valid    bool    // not read in or saved, but used by this App
	linkType string  // not read in or saved, but used by this App
}

// RelatedArticleLinks are related articles
type RelatedArticleLinks struct {
	URI      *string `bson:"uri,omitempty"  json:"uri,omitempty"`
	valid    bool    // not read in or saved, but used by this App
	linkType string  // not read in or saved, but used by this App
}

// RelatedMethodLinks are related methodologies
type RelatedMethodLinks struct {
	URI      *string `bson:"uri,omitempty"  json:"uri,omitempty"`
	valid    bool    // not read in or saved, but used by this App
	linkType string  // not read in or saved, but used by this App
}

// RelatedMethodologyArticleLinks are related methodology articles
type RelatedMethodologyArticleLinks struct {
	URI      *string `bson:"uri,omitempty"  json:"uri,omitempty"`
	valid    bool    // not read in or saved, but used by this App
	linkType string  // not read in or saved, but used by this App
}

// HighlightedContentLinks are highlighted content
type HighlightedContentLinks struct {
	URI      *string `bson:"uri,omitempty"  json:"uri,omitempty"`
	valid    bool    // not read in or saved, but used by this App
	linkType string  // not read in or saved, but used by this App
}

// ContactDetails represents an object containing information of the contact
type ContactDetails struct {
	Email     *string `bson:"email,omitempty"      json:"email,omitempty"`
	Name      *string `bson:"name,omitempty"       json:"name,omitempty"`
	Telephone *string `bson:"telephone,omitempty"  json:"telephone,omitempty"`
}

type allowedPageType int

const (
	pageBroken allowedPageType = iota + 1
	pageTopic
	pageTopicBroken
	pageTopicHighlightedLinks // Topic spotlight
	pageTopicSubtopicID
	pageTopicAndContent
	pageContent
	pageContentItems                     // Content Timeseries
	pageContentDatasets                  // Content Static datasets
	pageContentStatsBulletins            // Content Bulletins
	pageContentRelatedArticles           // Content Articles
	pageContentRelatedMethodology        // Content Methodologies
	pageContentRelatedMethodologyArticle // Content Methodology_articles
	pageContentHighlightedContent        // Content Spotlight
)

func pageTypeString(pType allowedPageType) string {
	var res string
	switch pType {
	case pageBroken:
		res = "pageBroken"
	case pageTopic:
		res = "sections"
	case pageTopicBroken:
		res = "pageTopicBroken"
	case pageTopicHighlightedLinks: // Topic spotlight
		res = "highlightedLinks"
	case pageTopicSubtopicID:
		res = "subtopicID"
	case pageContent:
		res = "pageContent"
	case pageContentItems: // Content Timeseries
		res = "items"
	case pageContentDatasets: // Content Static datasets
		res = "datasets"
	case pageContentStatsBulletins: // Content Bulletins
		res = "statsBulletins"
	case pageContentRelatedArticles: // Content Articles
		res = "relatedArticles"
	case pageContentRelatedMethodology: // Content Methodologies
		res = "relatedMethodology"
	case pageContentRelatedMethodologyArticle: // Content Methodology_articles
		res = "relatedMethodologyArticles"
	case pageContentHighlightedContent:
		res = "highlightedContent"
	default:
		res = "ERROR bad page type"
	}
	return res
}

// ===

/*
the output page types are derived from particular original page types as per:

	nodeHighlightedLinks               <== article, bulletin, compendium_landing_page
	contentItems                       <== timeseries
	contentDatasets                    <== compendium_data, dataset_landing_page
	contentStatsBulletins              <== bulletin
	contentRelatedArticles             <== article, article_download, compendium_landing_page
	contentRelatedMethodology          <== static_methodology, static_qmi
	contentRelatedMethodologyArticle   <== static_methodology, static_methodology_download
	contentHighlightedContent          <== article, bulletin, timeseries
*/

// ===
// pageShape used to determine how to decode a page
type pageShape struct {
	Type *string `bson:"type,omitempty"  json:"type,omitempty"`
}

// ===
type articleResponse struct {
	RelatedArticles           *[]relatedArticles           `bson:"relatedArticles,omitempty"            json:"relatedArticles,omitempty"`
	PdfTable                  *[]pdfTable                  `bson:"pdfTable,omitempty"                   json:"pdfTable,omitempty"`
	IsPrototypeArticle        *bool                        `bson:"isPrototypeArticle,omitempty"         json:"isPrototypeArticle,omitempty"`
	IsReleaseDateEnabled      *bool                        `bson:"isReleaseDateEnabled,omitempty"       json:"isReleaseDateEnabled,omitempty"`
	ImageURI                  *string                      `bson:"imageUri,omitempty"                   json:"imageUri,omitempty"`
	Sections                  *[]sections                  `bson:"sections,omitempty"                   json:"sections,omitempty"`
	Accordion                 *[]accordion                 `bson:"accordion,omitempty"                  json:"accordion,omitempty"`
	RelatedData               *[]relatedData               `bson:"relatedData,omitempty"                json:"relatedData,omitempty"`
	RelatedDocuments          *[]relatedDocuments          `bson:"relatedDocuments,omitempty"           json:"relatedDocuments,omitempty"`
	Charts                    *[]charts                    `bson:"charts,omitempty"                     json:"charts,omitempty"`
	Tables                    *[]tables                    `bson:"tables,omitempty"                     json:"tables,omitempty"`
	Images                    *[]images                    `bson:"images,omitempty"                     json:"images,omitempty"`
	Equations                 *[]equations                 `bson:"equations,omitempty"                  json:"equations,omitempty"`
	Links                     *[]links                     `bson:"links,omitempty"                      json:"links,omitempty"`
	Alerts                    *[]alerts                    `bson:"alerts,omitempty"                     json:"alerts,omitempty"`
	RelatedMethodology        *[]relatedMethodology        `bson:"relatedMethodology,omitempty"         json:"relatedMethodology,omitempty"`
	RelatedMethodologyArticle *[]relatedMethodologyArticle `bson:"relatedMethodologyArticle,omitempty"  json:"relatedMethodologyArticle,omitempty"`
	Versions                  *[]versions                  `bson:"versions,omitempty"                   json:"versions,omitempty"`
	Type                      *string                      `bson:"type,omitempty"                       json:"type,omitempty"`
	URI                       *string                      `bson:"uri,omitempty"                        json:"uri,omitempty"`
	Description               *cDescription                `bson:"description,omitempty"                json:"description,omitempty"`
	Topics                    *[]ctopics                   `bson:"topics,omitempty"                     json:"topics,omitempty"`
}

// ===
type articleDownloadResponse struct {
	Downloads                 *[]downloads                 `bson:"downloads,omitempty"                  json:"downloads,omitempty"`
	Markdown                  *[]string                    `bson:"markdown,omitempty"                   json:"markdown,omitempty"`
	RelatedData               *[]relatedData               `bson:"relatedData,omitempty"                json:"relatedData,omitempty"`
	RelatedDocuments          *[]relatedDocuments          `bson:"relatedDocuments,omitempty"           json:"relatedDocuments,omitempty"`
	Charts                    *[]charts                    `bson:"charts,omitempty"                     json:"charts,omitempty"`
	Tables                    *[]tables                    `bson:"tables,omitempty"                     json:"tables,omitempty"`
	Images                    *[]images                    `bson:"images,omitempty"                     json:"images,omitempty"`
	Equations                 *[]equations                 `bson:"equations,omitempty"                  json:"equations,omitempty"`
	Links                     *[]links                     `bson:"links,omitempty"                      json:"links,omitempty"`
	Alerts                    *[]alerts                    `bson:"alerts,omitempty"                     json:"alerts,omitempty"`
	RelatedMethodology        *[]relatedMethodology        `bson:"relatedMethodology,omitempty"         json:"relatedMethodology,omitempty"`
	RelatedMethodologyArticle *[]relatedMethodologyArticle `bson:"relatedMethodologyArticle,omitempty"  json:"relatedMethodologyArticle,omitempty"`
	Versions                  *[]versions                  `bson:"versions,omitempty"                   json:"versions,omitempty"`
	Type                      *string                      `bson:"type,omitempty"                       json:"type,omitempty"`
	URI                       *string                      `bson:"uri,omitempty"                        json:"uri,omitempty"`
	Description               *cDescription                `bson:"description,omitempty"                json:"description,omitempty"`
	Topics                    *[]ctopics                   `bson:"topics,omitempty"                     json:"topics,omitempty"`
}

// ===
type bulletinResponse struct {
	RelatedBulletins          *[]relatedBulletins          `bson:"relatedBulletins,omitempty"           json:"relatedBulletins,omitempty"`
	PdfTable                  *[]pdfTable                  `bson:"pdfTable,omitempty"                   json:"pdfTable,omitempty"`
	Sections                  *[]sections                  `bson:"sections,omitempty"                   json:"sections,omitempty"`
	Accordion                 *[]accordion                 `bson:"accordion,omitempty"                  json:"accordion,omitempty"`
	RelatedData               *[]relatedData               `bson:"relatedData,omitempty"                json:"relatedData,omitempty"`
	RelatedDocuments          *[]relatedDocuments          `bson:"relatedDocuments,omitempty"           json:"relatedDocuments,omitempty"`
	Charts                    *[]charts                    `bson:"charts,omitempty"                     json:"charts,omitempty"`
	Tables                    *[]tables                    `bson:"tables,omitempty"                     json:"tables,omitempty"`
	Images                    *[]images                    `bson:"images,omitempty"                     json:"images,omitempty"`
	Equations                 *[]equations                 `bson:"equations,omitempty"                  json:"equations,omitempty"`
	Links                     *[]links                     `bson:"links,omitempty"                      json:"links,omitempty"`
	Alerts                    *[]alerts                    `bson:"alerts,omitempty"                     json:"alerts,omitempty"`
	RelatedMethodology        *[]relatedMethodology        `bson:"relatedMethodology,omitempty"         json:"relatedMethodology,omitempty"`
	RelatedMethodologyArticle *[]relatedMethodologyArticle `bson:"relatedMethodologyArticle,omitempty"  json:"relatedMethodologyArticle,omitempty"`
	Versions                  *[]versions                  `bson:"versions,omitempty"                   json:"versions,omitempty"`
	Type                      *string                      `bson:"type,omitempty"                       json:"type,omitempty"`
	URI                       *string                      `bson:"uri,omitempty"                        json:"uri,omitempty"`
	Description               *cDescription                `bson:"description,omitempty"                json:"description,omitempty"`
	Topics                    *[]ctopics                   `bson:"topics,omitempty"                     json:"topics,omitempty"`
}

// ===
type compendiumDataResponse struct {
	Downloads                 *[]downloads                 `bson:"downloads,omitempty"                  json:"downloads,omitempty"`
	RelatedDatasets           *[]relatedDatasets           `bson:"relatedDatasets,omitempty"            json:"relatedDatasets,omitempty"`
	RelatedDocuments          *[]relatedDocuments          `bson:"relatedDocuments,omitempty"           json:"relatedDocuments,omitempty"`
	RelatedMethodology        *[]relatedMethodology        `bson:"relatedMethodology,omitempty"         json:"relatedMethodology,omitempty"`
	RelatedMethodologyArticle *[]relatedMethodologyArticle `bson:"relatedMethodologyArticle,omitempty"  json:"relatedMethodologyArticle,omitempty"`
	Versions                  *[]versions                  `bson:"versions,omitempty"                   json:"versions,omitempty"`
	Type                      *string                      `bson:"type,omitempty"                       json:"type,omitempty"`
	URI                       *string                      `bson:"uri,omitempty"                        json:"uri,omitempty"`
	Description               *cDescription                `bson:"description,omitempty"                json:"description,omitempty"`
	Topics                    *[]ctopics                   `bson:"topics,omitempty"                     json:"topics,omitempty"`
}

type compendiumLandingPageResponse struct {
	Datasets                  *[]datasets                  `bson:"datasets,omitempty"                   json:"datasets,omitempty"`
	Chapters                  *[]chapters                  `bson:"chapters,omitempty"                   json:"chapters,omitempty"`
	RelatedDocuments          *[]relatedDocuments          `bson:"relatedDocuments,omitempty"           json:"relatedDocuments,omitempty"`
	RelatedData               *[]relatedData               `bson:"relatedData,omitempty"                json:"relatedData,omitempty"`
	RelatedMethodology        *[]relatedMethodology        `bson:"relatedMethodology,omitempty"         json:"relatedMethodology,omitempty"`
	RelatedMethodologyArticle *[]relatedMethodologyArticle `bson:"relatedMethodologyArticle,omitempty"  json:"relatedMethodologyArticle,omitempty"`
	Alerts                    *[]alerts                    `bson:"alerts,omitempty"                     json:"alerts,omitempty"`
	Type                      *string                      `bson:"type,omitempty"                       json:"type,omitempty"`
	URI                       *string                      `bson:"uri,omitempty"                        json:"uri,omitempty"`
	Description               *cDescription                `bson:"description,omitempty"                json:"description,omitempty"`
	Topics                    *[]ctopics                   `bson:"topics,omitempty"                     json:"topics,omitempty"`
}

type datasetLandingPageResponse struct {
	Section                   *section                     `bson:"section,omitempty"                    json:"section,omitempty"`
	RelatedFilterableDatasets *[]relatedFilterableDatasets `bson:"relatedFilterableDatasets,omitempty"  json:"relatedFilterableDatasets,omitempty"`
	RelatedDatasets           *[]relatedDatasets           `bson:"relatedDatasets,omitempty"            json:"relatedDatasets,omitempty"`
	RelatedDocuments          *[]relatedDocuments          `bson:"relatedDocuments,omitempty"           json:"relatedDocuments,omitempty"`
	Datasets                  *[]datasets                  `bson:"datasets,omitempty"                   json:"datasets,omitempty"`
	Timeseries                *bool                        `bson:"timeseries,omitempty"                 json:"timeseries,omitempty"`
	Links                     *[]links                     `bson:"links,omitempty"                      json:"links,omitempty"`
	Alerts                    *[]alerts                    `bson:"alerts,omitempty"                     json:"alerts,omitempty"`
	RelatedMethodology        *[]relatedMethodology        `bson:"relatedMethodology,omitempty"         json:"relatedMethodology,omitempty"`
	RelatedMethodologyArticle *[]relatedMethodologyArticle `bson:"relatedMethodologyArticle,omitempty"  json:"relatedMethodologyArticle,omitempty"`
	Versions                  *[]versions                  `bson:"versions,omitempty"                   json:"versions,omitempty"`
	Type                      *string                      `bson:"type,omitempty"                       json:"type,omitempty"`
	URI                       *string                      `bson:"uri,omitempty"                        json:"uri,omitempty"`
	Description               *cDescription                `bson:"description,omitempty"                json:"description,omitempty"`
	Topics                    *[]ctopics                   `bson:"topics,omitempty"                     json:"topics,omitempty"`
}

type staticMethodologyResponse struct {
	PdfTable                  *[]pdfTable                  `bson:"pdfTable,omitempty"                   json:"pdfTable,omitempty"`
	RelatedDocuments          *[]relatedDocuments          `bson:"relatedDocuments,omitempty"           json:"relatedDocuments,omitempty"`
	Links                     *[]links                     `bson:"links,omitempty"                      json:"links,omitempty"`
	Downloads                 *[]downloads                 `bson:"downloads,omitempty"                  json:"downloads,omitempty"`
	Sections                  *[]sections                  `bson:"sections,omitempty"                   json:"sections,omitempty"`
	Accordion                 *[]accordion                 `bson:"accordion,omitempty"                  json:"accordion,omitempty"`
	RelatedData               *[]relatedData               `bson:"relatedData,omitempty"                json:"relatedData,omitempty"`
	Charts                    *[]charts                    `bson:"charts,omitempty"                     json:"charts,omitempty"`
	Tables                    *[]tables                    `bson:"tables,omitempty"                     json:"tables,omitempty"`
	Images                    *[]images                    `bson:"images,omitempty"                     json:"images,omitempty"`
	Equations                 *[]equations                 `bson:"equations,omitempty"                  json:"equations,omitempty"`
	Alerts                    *[]alerts                    `bson:"alerts,omitempty"                     json:"alerts,omitempty"`
	RelatedMethodology        *[]relatedMethodology        `bson:"relatedMethodology,omitempty"         json:"relatedMethodology,omitempty"`
	RelatedMethodologyArticle *[]relatedMethodologyArticle `bson:"relatedMethodologyArticle,omitempty"  json:"relatedMethodologyArticle,omitempty"`
	Versions                  *[]versions                  `bson:"versions,omitempty"                   json:"versions,omitempty"`
	Type                      *string                      `bson:"type,omitempty"                       json:"type,omitempty"`
	URI                       *string                      `bson:"uri,omitempty"                        json:"uri,omitempty"`
	Description               *cDescription                `bson:"description,omitempty"                json:"description,omitempty"`
	Topics                    *[]ctopics                   `bson:"topics,omitempty"                     json:"topics,omitempty"`
}

type staticMethodologyDownloadResponse struct {
	RelatedDocuments *[]relatedDocuments `bson:"relatedDocuments,omitempty"  json:"relatedDocuments,omitempty"`
	RelatedDatasets  *[]relatedDatasets  `bson:"relatedDatasets,omitempty"   json:"relatedDatasets,omitempty"`
	PdfDownloads     *[]pdfDownloads     `bson:"pdfDownloads,omitempty"      json:"pdfDownloads,omitempty"`
	Alerts           *[]alerts           `bson:"alerts,omitempty"            json:"alerts,omitempty"`
	Downloads        *[]downloads        `bson:"downloads,omitempty"         json:"downloads,omitempty"`
	Markdown         *[]string           `bson:"markdown,omitempty"          json:"markdown,omitempty"`
	Links            *[]links            `bson:"links,omitempty"             json:"links,omitempty"`
	Type             *string             `bson:"type,omitempty"              json:"type,omitempty"`
	URI              *string             `bson:"uri,omitempty"               json:"uri,omitempty"`
	Description      *cDescription       `bson:"description,omitempty"       json:"description,omitempty"`
	Topics           *[]ctopics          `bson:"topics,omitempty"            json:"topics,omitempty"`
}

type staticQmiResponse struct {
	RelatedDocuments *[]relatedDocuments `bson:"relatedDocuments,omitempty"  json:"relatedDocuments,omitempty"`
	RelatedDatasets  *[]relatedDatasets  `bson:"relatedDatasets,omitempty"   json:"relatedDatasets,omitempty"`
	Downloads        *[]downloads        `bson:"downloads,omitempty"         json:"downloads,omitempty"`
	Markdown         *[]string           `bson:"markdown,omitempty"          json:"markdown,omitempty"`
	Links            *[]links            `bson:"links,omitempty"             json:"links,omitempty"`
	Type             *string             `bson:"type,omitempty"              json:"type,omitempty"`
	URI              *string             `bson:"uri,omitempty"               json:"uri,omitempty"`
	Description      *cDescription       `bson:"description,omitempty"       json:"description,omitempty"`
	Topics           *[]ctopics          `bson:"topics,omitempty"            json:"topics,omitempty"`
}

type timeseriesResponse struct {
	Years            *[]years            `bson:"years,omitempty"             json:"years,omitempty"`
	Quarters         *[]quarters         `bson:"quarters,omitempty"          json:"quarters,omitempty"`
	Months           *[]months           `bson:"months,omitempty"            json:"months,omitempty"`
	SourceDatasets   *[]sourceDatasets   `bson:"sourceDatasets,omitempty"    json:"sourceDatasets,omitempty"`
	RelatedDatasets  *[]relatedDatasets  `bson:"relatedDatasets,omitempty"   json:"relatedDatasets,omitempty"`
	Section          *section            `bson:"section,omitempty"           json:"section,omitempty"`
	Notes            *[]string           `bson:"notes,omitempty"             json:"notes,omitempty"`
	RelatedDocuments *[]relatedDocuments `bson:"relatedDocuments,omitempty"  json:"relatedDocuments,omitempty"`
	RelatedData      *[]relatedData      `bson:"relatedData,omitempty"       json:"relatedData,omitempty"`
	Alerts           *[]alerts           `bson:"alerts,omitempty"            json:"alerts,omitempty"`
	Versions         *[]versions         `bson:"versions,omitempty"          json:"versions,omitempty"`
	Type             *string             `bson:"type,omitempty"              json:"type,omitempty"`
	URI              *string             `bson:"uri,omitempty"               json:"uri,omitempty"`
	Description      *cDescription       `bson:"description,omitempty"       json:"description,omitempty"`
	Topics           *[]ctopics          `bson:"topics,omitempty"            json:"topics,omitempty"`
}

type chartResponse struct {
	Subtitle           *string `bson:"subtitle,omitempty"            json:"subtitle,omitempty"`
	Filename           *string `bson:"filename,omitempty"            json:"filename,omitempty"`
	Source             *string `bson:"source,omitempty"              json:"source,omitempty"`
	Notes              *string `bson:"notes,omitempty"               json:"notes,omitempty"`
	AltText            *string `bson:"altText,omitempty"             json:"altText,omitempty"`
	LabelInterval      *string `bson:"labelInterval,omitempty"       json:"labelInterval,omitempty"`
	DecimalPlaces      *string `bson:"decimalPlaces,omitempty"       json:"decimalPlaces,omitempty"`
	DecimalPlacesYaxis *string `bson:"decimalPlacesYaxis,omitempty"  json:"decimalPlacesYaxis,omitempty"`
	Palette            *string `bson:"palette,omitempty"             json:"palette,omitempty"`
	XAxisPos           *string `bson:"xAxisPos,omitempty"            json:"xAxisPos,omitempty"`
	YAxisPos           *string `bson:"yAxisPos,omitempty"            json:"yAxisPos,omitempty"`
	YAxisMax           *string `bson:"yAxisMax,omitempty"            json:"yAxisMax,omitempty"`
	YMin               *string `bson:"yMin,omitempty"                json:"yMin,omitempty"`
	YMax               *string `bson:"yMax,omitempty"                json:"yMax,omitempty"`
	YxisInterval       *string `bson:"yAxisInterval,omitempty"       json:"yAxisInterval,omitempty"`
	Highlight          *string `bson:"highlight,omitempty"           json:"highlight,omitempty"`
	Alpha              *string `bson:"alpha,omitempty"               json:"alpha,omitempty"`
	Unit               *string `bson:"unit,omitempty"                json:"unit,omitempty"`
	XAxisLabel         *string `bson:"xAxisLabel,omitempty"          json:"xAxisLabel,omitempty"`
	AspectRatio        *string `bson:"aspectRatio,omitempty"         json:"aspectRatio,omitempty"`
	ChartType          *string `bson:"chartType,omitempty"           json:"chartType,omitempty"`
	// NOTE: 'go' does not maintain the order of the items in the 'map',
	//       though the order of the array is maintained in 'Data' below
	Data            *[]map[string]string `bson:"data,omitempty"             json:"data,omitempty"`
	Headers         *[]string            `bson:"headers,omitempty"          json:"headers,omitempty"`
	Series          *[]string            `bson:"series,omitempty"           json:"series,omitempty"`
	Categories      *[]string            `bson:"categories,omitempty"       json:"categories,omitempty"`
	ChartTypes      *map[string]string   `bson:"chartTypes,omitempty"       json:"chartTypes,omitempty"`
	LineTypes       *map[string]string   `bson:"lineTypes,omitempty"        json:"lineTypes,omitempty"`
	Groups          *[][]string          `bson:"groups,omitempty"           json:"groups,omitempty"`
	StartFromZero   *bool                `bson:"startFromZero,omitempty"    json:"startFromZero,omitempty"`
	FinishAtHundred *bool                `bson:"finishAtHundred,omitempty"  json:"finishAtHundred,omitempty"`
	IsStacked       *bool                `bson:"isStacked,omitempty"        json:"isStacked,omitempty"`
	IsReversed      *bool                `bson:"isReversed,omitempty"       json:"isReversed,omitempty"`
	ShowTooltip     *bool                `bson:"showTooltip,omitempty"      json:"showTooltip,omitempty"`
	ShowMarker      *bool                `bson:"showMarker,omitempty"       json:"showMarker,omitempty"`
	HasLineBreak    *bool                `bson:"hasLineBreak,omitempty"     json:"hasLineBreak,omitempty"`
	HasConnectNull  *bool                `bson:"hasConnectNull,omitempty"   json:"hasConnectNull,omitempty"`
	IsEditor        *bool                `bson:"isEditor,omitempty"         json:"isEditor,omitempty"`
	Annotations     *[]annotations       `bson:"annotations,omitempty"      json:"annotations,omitempty"`
	Devices         *devices             `bson:"devices,omitempty"          json:"devices,omitempty"`
	Files           *[]files             `bson:"files,omitempty"            json:"files,omitempty"`
	Title           *string              `bson:"title,omitempty"            json:"title,omitempty"`
	Type            *string              `bson:"type,omitempty"             json:"type,omitempty"`
	URI             *string              `bson:"uri,omitempty"              json:"uri,omitempty"`
}

type tableResponse struct {
	Filename       *string        `bson:"filename,omitempty"        json:"filename,omitempty"`
	FirstLineTitle *bool          `bson:"firstLineTitle,omitempty"  json:"firstLineTitle,omitempty"`
	HeaderRows     *string        `bson:"headerRows,omitempty"      json:"headerRows,omitempty"`
	Modifications  *modifications `bson:"modifications,omitempty"   json:"modifications,omitempty"`
	Files          *[]files       `bson:"files,omitempty"           json:"files,omitempty"`
	Title          *string        `bson:"title,omitempty"           json:"title,omitempty"`
	Type           *string        `bson:"type,omitempty"            json:"type,omitempty"`
	URI            *string        `bson:"uri,omitempty"             json:"uri,omitempty"`
}

type equationResponse struct {
	Filename *string  `bson:"filename,omitempty"  json:"filename,omitempty"`
	Content  *string  `bson:"content,omitempty"   json:"content,omitempty"`
	Files    *[]files `bson:"files,omitempty"     json:"files,omitempty"`
	Title    *string  `bson:"title,omitempty"     json:"title,omitempty"`
	Type     *string  `bson:"type,omitempty"      json:"type,omitempty"`
	URI      *string  `bson:"uri,omitempty"       json:"uri,omitempty"`
}

type imageResponse struct {
	Subtitle *string  `bson:"subtitle,omitempty"  json:"subtitle,omitempty"`
	Filename *string  `bson:"filename,omitempty"  json:"filename,omitempty"`
	Source   *string  `bson:"source,omitempty"    json:"source,omitempty"`
	Notes    *string  `bson:"notes,omitempty"     json:"notes,omitempty"`
	AltText  *string  `bson:"altText,omitempty"   json:"altText,omitempty"`
	Files    *[]files `bson:"files,omitempty"     json:"files,omitempty"`
	Title    *string  `bson:"title,omitempty"     json:"title,omitempty"`
	Type     *string  `bson:"type,omitempty"      json:"type,omitempty"`
	URI      *string  `bson:"uri,omitempty"       json:"uri,omitempty"`
}

type releaseResponse struct {
	Markdown                  *[]string                    `bson:"markdown,omitempty"                   json:"markdown,omitempty"`
	RelatedDocuments          *[]relatedDocuments          `bson:"relatedDocuments,omitempty"           json:"relatedDocuments,omitempty"`
	RelatedDatasets           *[]relatedDatasets           `bson:"relatedDatasets,omitempty"            json:"relatedDatasets,omitempty"`
	RelatedMethodology        *[]relatedMethodology        `bson:"relatedMethodology,omitempty"         json:"relatedMethodology,omitempty"`
	RelatedMethodologyArticle *[]relatedMethodologyArticle `bson:"relatedMethodologyArticle,omitempty"  json:"relatedMethodologyArticle,omitempty"`
	Links                     *[]links                     `bson:"links,omitempty"                      json:"links,omitempty"`
	DateChanges               *[]string                    `bson:"dateChanges,omitempty"                json:"dateChanges,omitempty"`
	Type                      *string                      `bson:"type,omitempty"                       json:"type,omitempty"`
	URI                       *string                      `bson:"uri,omitempty"                        json:"uri,omitempty"`
	Description               *cDescription                `bson:"description,omitempty"                json:"description,omitempty"`
}

type listResponse struct {
	Type     *string `bson:"type,omitempty"      json:"type,omitempty"`
	ListType *string `bson:"listType,omitempty"  json:"listType,omitempty"`
	URI      *string `bson:"uri,omitempty"       json:"uri,omitempty"`
	Result   *result `bson:"result,omitempty"    json:"result,omitempty"`
}

type result struct {
	NumberOfResults *int              `bson:"numberOfResults,omitempty"  json:"numberOfResults,omitempty"`
	Took            *int              `bson:"took,omitempty"             json:"took,omitempty"`
	Results         *[]results        `bson:"results,omitempty"          json:"results,omitempty"`
	Suggestions     *[]string         `bson:"suggestions,omitempty"      json:"suggestions,omitempty"`
	DocCounts       *map[string]int64 `bson:"docCounts,omitempty"        json:"docCounts,omitempty"`
	Paginator       *paginator        `bson:"paginator,omitempty"        json:"paginator,omitempty"`
	// NOTE: DocCounts has not been seen in any data read and its type has been derived from a chat with Jon who
	//       provided the following link to its java code from where the type now used has been derived:
	// https://github.com/ONSdigital/babbage/blob/fdfaa41528649c3dfb1165634330bfcb2c535fed/src/main/java/com/github/onsdigital/babbage/search/model/SearchResult.java#L43
	SortBy *string `bson:"sortBy,omitempty"           json:"sortBy,omitempty"`
}

type results struct {
	ResultType  *string            `bson:"_type,omitempty"        json:"_type,omitempty"`
	Description *resultDescription `bson:"description,omitempty"  json:"description,omitempty"`
	SearchBoost *[]string          `bson:"searchBoost,omitempty"  json:"searchBoost,omitempty"`
	Type        *string            `bson:"type,omitempty"         json:"type,omitempty"`
	URI         *string            `bson:"uri,omitempty"          json:"uri,omitempty"`
	Topics      *[]string          `bson:"topics,omitempty"       json:"topics,omitempty"`
	//	Topics      *[]resultsTopics   `bson:"topics,omitempty"       json:"topics,omitempty"`
	// for above, see link:
	// https://github.com/ONSdigital/zebedee/blob/957fb22f141546ae056f76bf1dbc23f4df8407de/zebedee-reader/src/main/java/com/github/onsdigital/zebedee/search/model/SearchDocument.java#L16
	// and then:
	// https://docs.oracle.com/javase/7/docs/api/java/net/URI.html
	// Jon suggests java struct: (but the above works ..)
	/*
		public URI(String scheme,
		   String userInfo,
		   String host,
		   int port,
		   String path,
		   String query,
		   String fragment)
	*/
}

/*
NOTE: this may need to be used in above commented out code at some point ..

type resultsTopics struct {
	UserInfo *string `bson:"userInfo,omitempty"  json:"userInfo,omitempty"`
	Host     *string `bson:"host,omitempty"      json:"host,omitempty"`
	Port     *int    `bson:"port,omitempty"      json:"port,omitempty"`
	Path     *string `bson:"path,omitempty"      json:"path,omitempty"`
	Query    *string `bson:"query,omitempty"     json:"query,omitempty"`
	Fragment *string `bson:"fragment,omitempty"  json:"fragment,omitempty"`
}*/

type paginator struct {
	NumberOfPages *int   `bson:"numberOfPages,omitempty"  json:"numberOfPages,omitempty"`
	CurrentPage   *int   `bson:"currentPage,omitempty"    json:"currentPage,omitempty"`
	Start         *int   `bson:"start,omitempty"          json:"start,omitempty"`
	End           *int   `bson:"end,omitempty"            json:"end,omitempty"`
	Pages         *[]int `bson:"pages,omitempty"          json:"pages,omitempty"`
}

type resultDescription struct {
	Summary           *string         `bson:"summary,omitempty"            json:"summary,omitempty"`
	NextRelease       *string         `bson:"nextRelease,omitempty"        json:"nextRelease,omitempty"`
	Keywords          *[]string       `bson:"keywords,omitempty"           json:"keywords,omitempty"`
	ReleaseDate       *string         `bson:"releaseDate,omitempty"        json:"releaseDate,omitempty"`
	Language          *string         `bson:"language,omitempty"           json:"language,omitempty"`
	DatasetID         *string         `bson:"datasetId,omitempty"          json:"datasetId,omitempty"`
	Edition           *string         `bson:"edition,omitempty"            json:"edition,omitempty"`
	Source            *string         `bson:"source,omitempty"             json:"source,omitempty"`
	Title             *string         `bson:"title,omitempty"              json:"title,omitempty"`
	MetaDescription   *string         `bson:"metaDescription,omitempty"    json:"metaDescription,omitempty"`
	NationalStatistic *bool           `bson:"nationalStatistic,omitempty"  json:"nationalStatistic,omitempty"`
	Headline1         *string         `bson:"headline1,omitempty"          json:"headline1,omitempty"`
	Headline3         *string         `bson:"headline3,omitempty"          json:"headline3,omitempty"`
	Headline2         *string         `bson:"headline2,omitempty"          json:"headline2,omitempty"`
	Abstract          *string         `bson:"_abstract,omitempty"          json:"_abstract,omitempty"`
	LatestRelease     *bool           `bson:"latestRelease,omitempty"      json:"latestRelease,omitempty"`
	Unit              *string         `bson:"unit,omitempty"               json:"unit,omitempty"`
	Contacts          *ContactDetails `bson:"contact,omitempty"            json:"contact,omitempty"`
	PreUnit           *string         `bson:"preUnit,omitempty"            json:"preUnit,omitempty"`
}

type staticPageResponse struct {
	Charts      *[]charts     `bson:"charts,omitempty"       json:"charts,omitempty"`
	Tables      *[]tables     `bson:"tables,omitempty"       json:"tables,omitempty"`
	Images      *[]images     `bson:"images,omitempty"       json:"images,omitempty"`
	Equations   *[]equations  `bson:"equations,omitempty"    json:"equations,omitempty"`
	Downloads   *[]downloads  `bson:"downloads,omitempty"    json:"downloads,omitempty"`
	Markdown    *[]string     `bson:"markdown,omitempty"     json:"markdown,omitempty"`
	Links       *[]links      `bson:"links,omitempty"        json:"links,omitempty"`
	Type        *string       `bson:"type,omitempty"         json:"type,omitempty"`
	URI         *string       `bson:"uri,omitempty"          json:"uri,omitempty"`
	Description *cDescription `bson:"description,omitempty"  json:"description,omitempty"`
}

type staticAdhocResponse struct {
	Downloads   *[]downloads  `bson:"downloads,omitempty"    json:"downloads,omitempty"`
	Markdown    *[]string     `bson:"markdown,omitempty"     json:"markdown,omitempty"`
	Links       *[]links      `bson:"links,omitempty"        json:"links,omitempty"`
	Type        *string       `bson:"type,omitempty"         json:"type,omitempty"`
	URI         *string       `bson:"uri,omitempty"          json:"uri,omitempty"`
	Description *cDescription `bson:"description,omitempty"  json:"description,omitempty"`
}

type referenceTablesResponse struct {
	Migrated           *bool                 `bson:"migrated,omitempty"            json:"migrated,omitempty"`
	Downloads          *[]downloads          `bson:"downloads,omitempty"           json:"downloads,omitempty"`
	RelatedDocuments   *[]relatedDocuments   `bson:"relatedDocuments,omitempty"    json:"relatedDocuments,omitempty"`
	RelatedMethodology *[]relatedMethodology `bson:"relatedMethodology,omitempty"  json:"relatedMethodology,omitempty"`
	Type               *string               `bson:"type,omitempty"                json:"type,omitempty"`
	URI                *string               `bson:"uri,omitempty"                 json:"uri,omitempty"`
	Description        *cDescription         `bson:"description,omitempty"         json:"description,omitempty"`
}

type compendiumChapterResponse struct {
	PdfTable                  *[]pdfTable                  `bson:"pdfTable,omitempty"                   json:"pdfTable,omitempty"`
	Sections                  *[]sections                  `bson:"sections,omitempty"                   json:"sections,omitempty"`
	Accordion                 *[]accordion                 `bson:"accordion,omitempty"                  json:"accordion,omitempty"`
	RelatedData               *[]relatedData               `bson:"relatedData,omitempty"                json:"relatedData,omitempty"`
	RelatedDocuments          *[]relatedDocuments          `bson:"relatedDocuments,omitempty"           json:"relatedDocuments,omitempty"`
	Charts                    *[]charts                    `bson:"charts,omitempty"                     json:"charts,omitempty"`
	Tables                    *[]tables                    `bson:"tables,omitempty"                     json:"tables,omitempty"`
	Images                    *[]images                    `bson:"images,omitempty"                     json:"images,omitempty"`
	Equations                 *[]equations                 `bson:"equations,omitempty"                  json:"equations,omitempty"`
	Links                     *[]links                     `bson:"links,omitempty"                      json:"links,omitempty"`
	Alerts                    *[]alerts                    `bson:"alerts,omitempty"                     json:"alerts,omitempty"`
	RelatedMethodology        *[]relatedMethodology        `bson:"relatedMethodology,omitempty"         json:"relatedMethodology,omitempty"`
	RelatedMethodologyArticle *[]relatedMethodologyArticle `bson:"relatedMethodologyArticle,omitempty"  json:"relatedMethodologyArticle,omitempty"`
	Versions                  *[]versions                  `bson:"versions,omitempty"                   json:"versions,omitempty"`
	Type                      *string                      `bson:"type,omitempty"                       json:"type,omitempty"`
	URI                       *string                      `bson:"uri,omitempty"                        json:"uri,omitempty"`
	Description               *cDescription                `bson:"description,omitempty"                json:"description,omitempty"`
}

type staticLandingPageResponse struct {
	Sections    *[]sections   `bson:"sections,omitempty"     json:"sections,omitempty"`
	Markdown    *[]string     `bson:"markdown,omitempty"     json:"markdown,omitempty"`
	Links       *[]links      `bson:"links,omitempty"        json:"links,omitempty"`
	Type        *string       `bson:"type,omitempty"         json:"type,omitempty"`
	URI         *string       `bson:"uri,omitempty"          json:"uri,omitempty"`
	Description *cDescription `bson:"description,omitempty"  json:"description,omitempty"`
}

type staticArticleResponse struct {
	Links       *[]links      `bson:"links,omitempty"        json:"links,omitempty"`
	Downloads   *[]downloads  `bson:"downloads,omitempty"    json:"downloads,omitempty"`
	Sections    *[]sections   `bson:"sections,omitempty"     json:"sections,omitempty"`
	Accordion   *[]accordion  `bson:"accordion,omitempty"    json:"accordion,omitempty"`
	Charts      *[]charts     `bson:"charts,omitempty"       json:"charts,omitempty"`
	Tables      *[]tables     `bson:"tables,omitempty"       json:"tables,omitempty"`
	Images      *[]images     `bson:"images,omitempty"       json:"images,omitempty"`
	Equations   *[]equations  `bson:"equations,omitempty"    json:"equations,omitempty"`
	Alerts      *[]alerts     `bson:"alerts,omitempty"       json:"alerts,omitempty"`
	Type        *string       `bson:"type,omitempty"         json:"type,omitempty"`
	URI         *string       `bson:"uri,omitempty"          json:"uri,omitempty"`
	Description *cDescription `bson:"description,omitempty"  json:"description,omitempty"`
}

type datasetResponse struct {
	Downloads          *[]downloads          `bson:"downloads,omitempty"           json:"downloads,omitempty"`
	SupplementaryFiles *[]supplementaryFiles `bson:"supplementaryFiles,omitempty"  json:"supplementaryFiles,omitempty"`
	Versions           *[]versions           `bson:"versions,omitempty"            json:"versions,omitempty"`
	Type               *string               `bson:"type,omitempty"                json:"type,omitempty"`
	URI                *string               `bson:"uri,omitempty"                 json:"uri,omitempty"`
	Description        *cDescription         `bson:"description,omitempty"         json:"description,omitempty"`
}

type timeseriesDatasetResponse struct {
	Downloads          *[]downloads          `bson:"downloads,omitempty"           json:"downloads,omitempty"`
	SupplementaryFiles *[]supplementaryFiles `bson:"supplementaryFiles,omitempty"  json:"supplementaryFiles,omitempty"`
	Versions           *[]versions           `bson:"versions,omitempty"            json:"versions,omitempty"`
	Type               *string               `bson:"type,omitempty"                json:"type,omitempty"`
	URI                *string               `bson:"uri,omitempty"                 json:"uri,omitempty"`
	Description        *cDescription         `bson:"description,omitempty"         json:"description,omitempty"`
}

type supplementaryFiles struct {
	Title *string `bson:"title,omitempty"  json:"title,omitempty"`
	File  *string `bson:"file,omitempty"   json:"file,omitempty"`
}

type modifications struct {
	RowsExcluded  *[]int `bson:"rowsExcluded,omitempty"   json:"rowsExcluded,omitempty"`
	HeaderRows    *[]int `bson:"headerRows,omitempty"     json:"headerRows,omitempty"`
	HeaderColumns *[]int `bson:"headerColumns,omitempty"  json:"headerColumns,omitempty"`
}

type files struct {
	Type     *string `bson:"type,omitempty"      json:"type,omitempty"`
	Filename *string `bson:"filename,omitempty"  json:"filename,omitempty"`
	FileType *string `bson:"fileType,omitempty"  json:"fileType,omitempty"`
}

type annotations struct {
	ID          *string  `bson:"id,omitempty"           json:"id,omitempty"`
	X           *string  `bson:"x,omitempty"            json:"x,omitempty"`
	Y           *string  `bson:"y,omitempty"            json:"y,omitempty"`
	Title       *string  `bson:"title,omitempty"        json:"title,omitempty"`
	Orientation *string  `bson:"orientation,omitempty"  json:"orientation,omitempty"`
	IsHidden    *bool    `bson:"isHidden,omitempty"     json:"isHidden,omitempty"`
	IsPlotline  *bool    `bson:"isPlotline,omitempty"   json:"isPlotline,omitempty"`
	IsPlotband  *bool    `bson:"isPlotband,omitempty"   json:"isPlotband,omitempty"`
	BandWidth   *string  `bson:"bandWidth,omitempty"    json:"bandWidth,omitempty"`
	Width       *string  `bson:"width,omitempty"        json:"width,omitempty"`
	Height      *string  `bson:"height,omitempty"       json:"height,omitempty"`
	Devices     *devices `bson:"devices,omitempty"      json:"devices,omitempty"`
}

type devices struct {
	Sm *deviceType `bson:"sm,omitempty"  json:"sm,omitempty"`
	Md *deviceType `bson:"md,omitempty"  json:"md,omitempty"`
	Lg *deviceType `bson:"lg,omitempty"  json:"lg,omitempty"`
}

type deviceType struct {
	AspectRatio   *string `bson:"aspectRatio,omitempty"    json:"aspectRatio,omitempty"`
	LabelInterval *string `bson:"labelInterval,omitempty"  json:"labelInterval,omitempty"`
	IsHidden      *bool   `bson:"isHidden,omitempty"       json:"isHidden,omitempty"`
	XAxisLabel    *string `bson:"x,omitempty"              json:"x,omitempty"`
	Y             *string `bson:"y,omitempty"              json:"y,omitempty"`
}

type years struct {
	Date          *string `bson:"date,omitempty"           json:"date,omitempty"`
	Value         *string `bson:"value,omitempty"          json:"value,omitempty"`
	Label         *string `bson:"label,omitempty"          json:"label,omitempty"`
	Year          *string `bson:"year,omitempty"           json:"year,omitempty"`
	Month         *string `bson:"month,omitempty"          json:"month,omitempty"`
	Quarter       *string `bson:"quarter,omitempty"        json:"quarter,omitempty"`
	SourceDataset *string `bson:"sourceDataset,omitempty"  json:"sourceDataset,omitempty"`
	UpdateDate    *string `bson:"updateDate,omitempty"     json:"updateDate,omitempty"`
}

type quarters struct {
	Date          *string `bson:"date,omitempty"           json:"date,omitempty"`
	Value         *string `bson:"value,omitempty"          json:"value,omitempty"`
	Label         *string `bson:"label,omitempty"          json:"label,omitempty"`
	Year          *string `bson:"year,omitempty"           json:"year,omitempty"`
	Month         *string `bson:"month,omitempty"          json:"month,omitempty"`
	Quarter       *string `bson:"quarter,omitempty"        json:"quarter,omitempty"`
	SourceDataset *string `bson:"sourceDataset,omitempty"  json:"sourceDataset,omitempty"`
	UpdateDate    *string `bson:"updateDate,omitempty"     json:"updateDate,omitempty"`
}

type months struct {
	Date          *string `bson:"date,omitempty"           json:"date,omitempty"`
	Value         *string `bson:"value,omitempty"          json:"value,omitempty"`
	Label         *string `bson:"label,omitempty"          json:"label,omitempty"`
	Year          *string `bson:"year,omitempty"           json:"year,omitempty"`
	Month         *string `bson:"month,omitempty"          json:"month,omitempty"`
	Quarter       *string `bson:"quarter,omitempty"        json:"quarter,omitempty"`
	SourceDataset *string `bson:"sourceDataset,omitempty"  json:"sourceDataset,omitempty"`
	UpdateDate    *string `bson:"updateDate,omitempty"     json:"updateDate,omitempty"`
}
type downloads struct {
	Title           *string `bson:"title,omitempty"            json:"title,omitempty"`
	File            *string `bson:"file,omitempty"             json:"file,omitempty"`
	FileDescription *string `bson:"fileDescription,omitempty"  json:"fileDescription,omitempty"`
}
type relatedArticles struct {
	Article *string `bson:"article,omitempty"  json:"article,omitempty"`
}

type relatedBulletins struct {
	Bulletin *string `bson:"bulletin,omitempty"  json:"bulletin,omitempty"`
}

type pdfTable struct {
	Title *string `bson:"title,omitempty"  json:"title,omitempty"`
	File  *string `bson:"file,omitempty"   json:"file,omitempty"`
}

type pdfDownloads struct {
	Title *string `bson:"title,omitempty"  json:"title,omitempty"`
	File  *string `bson:"file,omitempty"   json:"file,omitempty"`
}

type sections struct {
	Summary  *string `bson:"summary,omitempty"   json:"summary,omitempty"`
	Title    *string `bson:"title,omitempty"     json:"title,omitempty"`
	URI      *string `bson:"uri,omitempty"       json:"uri,omitempty"`
	Markdown *string `bson:"markdown,omitempty"  json:"markdown,omitempty"`
}

type accordion struct {
	Title    *string `bson:"title,omitempty"     json:"title,omitempty"`
	Markdown *string `bson:"markdown,omitempty"  json:"markdown,omitempty"`
}

type chapters struct {
	URI *string `bson:"uri,omitempty"  json:"uri,omitempty"`
}

type relatedData struct {
	URI *string `bson:"uri,omitempty"  json:"uri,omitempty"`
}

type section struct {
	Markdown *string `bson:"markdown,omitempty"  json:"markdown,omitempty"`
}

// relatedFilterableDatasets is the sub page links
type relatedFilterableDatasets struct {
	URI *string `bson:"uri,omitempty"  json:"uri,omitempty"`
}

// relatedDatasets is the sub page links
type relatedDatasets struct {
	Title *string `bson:"title,omitempty"  json:"title,omitempty"`
	URI   *string `bson:"uri,omitempty"    json:"uri,omitempty"`
}

type sourceDatasets struct {
	URI *string `bson:"uri,omitempty"  json:"uri,omitempty"`
}

// relatedDocuments is the sub page links
type relatedDocuments struct {
	Title *string `bson:"title,omitempty"  json:"title,omitempty"`
	URI   *string `bson:"uri,omitempty"    json:"uri,omitempty"`
}

// datasets is the sub page links
type datasets struct {
	URI *string `bson:"uri,omitempty"  json:"uri,omitempty"`
}

type charts struct {
	Title    *string `bson:"title,omitempty"     json:"title,omitempty"`
	Filename *string `bson:"filename,omitempty"  json:"filename,omitempty"`
	Version  *string `bson:"version,omitempty"   json:"version,omitempty"`
	URI      *string `bson:"uri,omitempty"       json:"uri,omitempty"`
}

type tables struct {
	Title    *string `bson:"title,omitempty"     json:"title,omitempty"`
	Filename *string `bson:"filename,omitempty"  json:"filename,omitempty"`
	Version  *string `bson:"version,omitempty"   json:"version,omitempty"`
	URI      *string `bson:"uri,omitempty"       json:"uri,omitempty"`
}

type images struct {
	Title    *string `bson:"title,omitempty"     json:"title,omitempty"`
	Filename *string `bson:"filename,omitempty"  json:"filename,omitempty"`
	URI      *string `bson:"uri,omitempty"       json:"uri,omitempty"`
}

type equations struct {
	Title    *string `bson:"title,omitempty"     json:"title,omitempty"`
	Filename *string `bson:"filename,omitempty"  json:"filename,omitempty"`
	URI      *string `bson:"uri,omitempty"       json:"uri,omitempty"`
}

// links is the sub page links
type links struct {
	Title *string `bson:"title,omitempty"  json:"title,omitempty"`
	URI   *string `bson:"uri,omitempty"    json:"uri,omitempty"`
}

// alerts is the sub page links
type alerts struct {
	Date     *string `bson:"date,omitempty"      json:"date,omitempty"`
	Markdown *string `bson:"markdown,omitempty"  json:"markdown,omitempty"`
	Type     *string `bson:"type,omitempty"      json:"type,omitempty"`
}

// relatedMethodology is the sub page links
type relatedMethodology struct {
	URI *string `bson:"uri,omitempty"  json:"uri,omitempty"`
}

// relatedMethodologyArticle is the sub page links
type relatedMethodologyArticle struct {
	Title *string `bson:"title,omitempty"  json:"title,omitempty"`
	URI   *string `bson:"uri,omitempty"    json:"uri,omitempty"`
}

type versions struct {
	URI              *string `bson:"uri,omitempty"               json:"uri,omitempty"`
	UpdateDate       *string `bson:"updateDate,omitempty"        json:"updateDate,omitempty"`
	CorrectionNotice *string `bson:"correctionNotice,omitempty"  json:"correctionNotice,omitempty"`
	Label            *string `bson:"label,omitempty"             json:"label,omitempty"`
}

// ctopics is the sub page links
type ctopics struct {
	URI *string `bson:"uri,omitempty"  json:"uri,omitempty"`
}

// cDescription is the page description
type cDescription struct {
	Finalised          *bool           `bson:"finalised,omitempty"           json:"finalised,omitempty"`
	Title              *string         `bson:"title,omitempty"               json:"title,omitempty"`
	Summary            *string         `bson:"summary,omitempty"             json:"summary,omitempty"`
	Keywords           *[]string       `bson:"keywords,omitempty"            json:"keywords,omitempty"`
	MetaDescription    *string         `bson:"metaDescription,omitempty"     json:"metaDescription,omitempty"`
	NationalStatistic  *bool           `bson:"nationalStatistic,omitempty"   json:"nationalStatistic,omitempty"`
	LatestRelease      *bool           `bson:"latestRelease,omitempty"       json:"latestRelease,omitempty"`
	Contacts           *ContactDetails `bson:"contact,omitempty"             json:"contact,omitempty"`
	ReleaseDate        *string         `bson:"releaseDate,omitempty"         json:"releaseDate,omitempty"`
	NextRelease        *string         `bson:"nextRelease,omitempty"         json:"nextRelease,omitempty"`
	Language           *string         `bson:"language,omitempty"            json:"language,omitempty"`
	Edition            *string         `bson:"edition,omitempty"             json:"edition,omitempty"`
	DatasetID          *string         `bson:"datasetId,omitempty"           json:"datasetId,omitempty"`
	DatasetURI         *string         `bson:"datasetUri,omitempty"          json:"datasetUri,omitempty"`
	Cdid               *string         `bson:"cdid,omitempty"                json:"cdid,omitempty"`
	Abstract           *string         `bson:"_abstract,omitempty"           json:"_abstract,omitempty"`
	Authors            *[]string       `bson:"authors,omitempty"             json:"authors,omitempty"`
	Headline           *string         `bson:"headline,omitempty"            json:"headline,omitempty"`
	Headline1          *string         `bson:"headline1,omitempty"           json:"headline1,omitempty"`
	Headline2          *string         `bson:"headline2,omitempty"           json:"headline2,omitempty"`
	Headline3          *string         `bson:"headline3,omitempty"           json:"headline3,omitempty"`
	Unit               *string         `bson:"unit,omitempty"                json:"unit,omitempty"`
	PreUnit            *string         `bson:"preUnit,omitempty"             json:"preUnit,omitempty"`
	Source             *string         `bson:"source,omitempty"              json:"source,omitempty"`
	Reference          *string         `bson:"reference,omitempty"           json:"reference,omitempty"`
	Cancelled          *bool           `bson:"cancelled,omitempty"           json:"cancelled,omitempty"`
	CancellationNotice *[]string       `bson:"cancellationNotice,omitempty"  json:"cancellationNotice,omitempty"`
	Published          *bool           `bson:"published,omitempty"           json:"published,omitempty"`
	ProvisionalDate    *string         `bson:"provisionalDate,omitempty"     json:"provisionalDate,omitempty"`
	MonthLabelStyle    *string         `bson:"monthLabelStyle,omitempty"     json:"monthLabelStyle,omitempty"`
	Date               *string         `bson:"date,omitempty"                json:"date,omitempty"`
	Number             *string         `bson:"number,omitempty"              json:"number,omitempty"`
	SurveyName         *string         `bson:"surveyName,omitempty"          json:"surveyName,omitempty"`
	Frequency          *string         `bson:"frequency,omitempty"           json:"frequency,omitempty"`
	Compilation        *string         `bson:"compilation,omitempty"         json:"compilation,omitempty"`
	GeographicCoverage *string         `bson:"geographicCoverage,omitempty"  json:"geographicCoverage,omitempty"`
	MetaCmd            *string         `bson:"metaCmd,omitempty"             json:"metaCmd,omitempty"`
	KeyNote            *string         `bson:"keyNote,omitempty"             json:"keyNote,omitempty"`
	SampleSize         *string         `bson:"sampleSize,omitempty"          json:"sampleSize,omitempty"`
	LastRevised        *string         `bson:"lastRevised,omitempty"         json:"lastRevised,omitempty"`
	VersionLabel       *string         `bson:"versionLabel,omitempty"        json:"versionLabel,omitempty"`
}

// ===

type contentInfo int

const (
	contentUnknown contentInfo = iota + 1
	contentNone
	contentExists
)

var pageCount int
var attemptedGetCount int

// The first use of this will increment it to 1 and save that.
// The value of '1' is reserved for the root node that is created using sub node
// index values that are determined as the code runs.
var indexNumber int = 0

var indexNames = make(map[int]string) // key: id, value: shortURI

type duplicateInfo struct {
	id        int
	pageType  allowedPageType
	parentURI string
	shortURI  string
}

var listOfDuplicateInfo []duplicateInfo

// store the shortUri and count
var appearanceInfo = make(map[string]int) // key: shortURI, value: count of pages it appears on

// create and store an index number for each page for use in creating mongo 'id' for each datafile in collection
// and for use in the cross referencing.
type pageData struct {
	id              int
	subSectionIndex int
	pageType        allowedPageType // the type of page that was trying to be read
	uriStatus       allowedPageType // the result of the page read: pageBroken, pageTopicBroken OR same as pageType indicating OK
	shortURI        string
	parentURI       string
	fixedPayload    []byte
	title           string
	description     string
}

var listOfPageData []pageData

func replaceUnicodeWithASCII(b []byte) []byte {
	l := len(b)
	var dst int

	for src := 0; src < l; src++ {
		if b[src] == '\\' {
			if src < l-5 {
				if b[src+1] == 'u' {
					// We have a unicode sequence, and we assume it is 6 characters long
					// get 4 hex characters
					hexstring := string([]byte{b[src+2], b[src+3], b[src+4], b[src+5]})
					num, err := hex.DecodeString(hexstring)
					if err != nil {
						// we can get an error whilst trying to decode an 'equation' type of page that has
						// special encodings to describe the equation drawn on the web page, a page like:
						// https://www.ons.gov.uk/employmentandlabourmarket/peopleinwork/earningsandworkinghours/articles/understandingthegenderpaygapintheuk/2018-01-17/3b240f04/data
						/*fmt.Printf("Error in DecodeString: %v\n", hexstring)
						fmt.Printf("byte 0: %v\n", b[src+2])
						fmt.Printf("byte 1: %v\n", b[src+3])
						fmt.Printf("byte 2: %v\n", b[src+4])
						fmt.Printf("byte 3: %v\n", b[src+5])
						panic(err)*/
						// So, we have to assume that this is an equation and just copy the character ..
						b[dst] = b[src]
					} else {
						b[dst] = num[1] // get ASCII character
						src += 5        // skip past unicode sequence - the for loop increment makes this an increase of 6
					}
				} else {
					b[dst] = b[src]
				}
				dst++
			}
		} else {
			b[dst] = b[src]
			dst++
		}
	}
	return b[0:dst]
}

// the following is examples of the \u0027  and its change into a single quote:
/*
 {"items":[{"uri":"/economy/economicoutputandproductivity/output/timeseries/k27q"},{"uri":"/economy/economicoutputandproductivity/output/timeseries/k222"},{"uri":"/economy/economicoutputandproductivity/output/timeseries/k27y"},{"uri":"/economy/economicoutputandproductivity/output/timeseries/k22a"}],"datasets":[{"uri":"/businessindustryandtrade/manufacturingandproductionindustry/datasets/topsiproductionandservicesturnover"},{"uri":"/businessindustryandtrade/manufacturingandproductionindustry/datasets/turnoverandordersintheproductionandservicesindustriesdataset"},{"uri":"/businessindustryandtrade/manufacturingandproductionindustry/datasets/topsimanufacturingexportturnover"},{"uri":"/businessindustryandtrade/manufacturingandproductionindustry/datasets/ukmanufacturerssalesbyproductprodcom"},{"uri":"/businessindustryandtrade/manufacturingandproductionindustry/datasets/ukmanufacturerssalesbyproductprodcomintermediateresults2013andfinalresults2012referencetables"}],"statsBulletins":[{"uri":"/businessindustryandtrade/manufacturingandproductionindustry/bulletins/ukmanufacturerssalesbyproductprodcom/latest"},{"uri":"/businessindustryandtrade/business/businessservices/bulletins/uknonfinancialbusinesseconomy/latest"}],"relatedArticles":[],"relatedMethodology":[{"uri":"/businessindustryandtrade/manufacturingandproductionindustry/qmis/ukmanufacturerssalesbyproductqmi"},{"uri":"/businessindustryandtrade/business/businessservices/qmis/monthlybusinesssurveyqmi"}],"relatedMethodologyArticle":[{"uri":"/businessindustryandtrade/manufacturingandproductionindustry/methodologies/ukmanufacturerssalesbyproductprodcom"}],"index":5,"type":"product_page","uri":"/businessindustryandtrade/manufacturingandproductionindustry","description":{"title":"Manufacturing and production industry","summary":"UK manufacturing and other production industries (such as mining and quarrying, energy supply, water supply and waste management), including total UK production output, and
UK manufactures\u0027
 sales by product and industrial division, with EU comparisons.","metaDescription":"UK manufacturing and other production industries (such as mining and quarrying, energy supply, water supply and waste management), including total UK production output, and UK manufactures\u0027 sales by product and industrial division, with EU comparisons.","unit":"","preUnit":"","source":""}}

{"items":[{"uri":"/economy/economicoutputandproductivity/output/timeseries/k27q"},{"uri":"/economy/economicoutputandproductivity/output/timeseries/k222"},{"uri":"/economy/economicoutputandproductivity/output/timeseries/k27y"},{"uri":"/economy/economicoutputandproductivity/output/timeseries/k22a"}],"datasets":[{"uri":"/businessindustryandtrade/manufacturingandproductionindustry/datasets/topsiproductionandservicesturnover"},{"uri":"/businessindustryandtrade/manufacturingandproductionindustry/datasets/turnoverandordersintheproductionandservicesindustriesdataset"},{"uri":"/businessindustryandtrade/manufacturingandproductionindustry/datasets/topsimanufacturingexportturnover"},{"uri":"/businessindustryandtrade/manufacturingandproductionindustry/datasets/ukmanufacturerssalesbyproductprodcom"},{"uri":"/businessindustryandtrade/manufacturingandproductionindustry/datasets/ukmanufacturerssalesbyproductprodcomintermediateresults2013andfinalresults2012referencetables"}],"statsBulletins":[{"uri":"/businessindustryandtrade/manufacturingandproductionindustry/bulletins/ukmanufacturerssalesbyproductprodcom/latest"},{"uri":"/businessindustryandtrade/business/businessservices/bulletins/uknonfinancialbusinesseconomy/latest"}],"relatedArticles":[],"relatedMethodology":[{"uri":"/businessindustryandtrade/manufacturingandproductionindustry/qmis/ukmanufacturerssalesbyproductqmi"},{"uri":"/businessindustryandtrade/business/businessservices/qmis/monthlybusinesssurveyqmi"}],"relatedMethodologyArticle":[{"uri":"/businessindustryandtrade/manufacturingandproductionindustry/methodologies/ukmanufacturerssalesbyproductprodcom"}],"index":5,"type":"product_page","uri":"/businessindustryandtrade/manufacturingandproductionindustry","description":{"title":"Manufacturing and production industry","summary":"UK manufacturing and other production industries (such as mining and quarrying, energy supply, water supply and waste management), including total UK production output, and
UK manufactures'
 sales by product and industrial division, with EU comparisons.","metaDescription":"UK manufacturing and other production industries (such as mining and quarrying, energy supply, water supply and waste management), including total UK production output, and UK manufactures' sales by product and industrial division, with EU comparisons.","unit":"","preUnit":"","source":""}}
*/

func doAndShowDelay71() {
	fmt.Printf("Got a 429, backing off for 71 seconds ..\n")
	for delay := 0; delay < 71; delay++ {
		time.Sleep(1 * time.Second)
		if delay%3 == 0 {
			fmt.Printf("Index: %d   Seconds remaining: %d\n", indexNumber, 71-delay)
		}
	}
}

// getPage returns pageURI, index, hasLink
func getPage(parentID int, graphVizFile io.Writer, parentURI, shortURI string) (string, int, bool, allowedPageType) {

	// Add prefix and '/data' to shortURI name
	//	fullURI := "https://www.ons.gov.uk" + shortURI + "/data"
	fullURI := "https://www.production.onsdigital.co.uk" + shortURI + "/data"

	if parentURI == "" && shortURI == rootPath {
		fullURI = rootURI + shortURI
	}

	// remove leading '/'
	gvPage := shortURI[1:]

	// formulate graphviz label
	gvPageLabel := "/" + strings.ReplaceAll(gvPage, "/", "\\n/")

	// replace remaining '/' with '_' to have a connection name that is suitable for graphviz
	gvPage = strings.ReplaceAll(gvPage, "/", "_")

	var response *http.Response
	var err error
	var attempts int

	for {
		attemptedGetCount++
		if cfg.PlayNice {
			// a little delay to play nice with ONS site and 'hopefully' not have cloudflare 'reset' the connection
			time.Sleep(100 * time.Millisecond)
		}
		response, err = http.Get(fullURI)
		if err != nil {
			fmt.Println(err)
			fmt.Printf("getPage: http.Get of fullURI failed\n")
			fmt.Printf("We now fabricate the response code to a 429 to instigate a retry after a delay 1\n")
			doAndShowDelay71()
		} else {
			if response.StatusCode == 429 {
				response.Body.Close()
				doAndShowDelay71()
			} else {
				break
			}
		}
		attempts++
		if attempts >= 3 {
			// Possible problems are:
			// 1. URI on ONS is broke
			// 2. ONS site is down
			// 3. Network connection to ONS is down
			// SO, give up on this URI ..

			// the node that is linked to is broken - that is it does not exist
			indexNumber++
			indexNames[indexNumber] = shortURI

			appearanceInfo[shortURI]++
			listOfDuplicateInfo = append(listOfDuplicateInfo, duplicateInfo{
				id:        indexNumber,
				pageType:  pageTopic,
				parentURI: parentURI,
				shortURI:  shortURI,
			})
			listOfPageData = append(listOfPageData, pageData{
				id:              indexNumber,
				subSectionIndex: parentID,
				pageType:        pageTopic,
				uriStatus:       pageTopicBroken,
				shortURI:        shortURI,
				parentURI:       "https://www.ons.gov.uk" + parentURI + "/data",
				fixedPayload:    []byte{},
			})

			fmt.Printf("\nERROR on ONS website /data field: %v\n\n", response.StatusCode)
			fmt.Printf("URI does not exist:  %v\n", fullURI)
			_, err := fmt.Fprintf(graphVizFile, "    subgraph %s {\n        %s [shape = box, color=red, style=bold, label = \"%s\"]\n    }\n", gvPage, gvPage, gvPageLabel+"\n ** MISSING: /data **"+fmt.Sprintf("\\n%v", indexNumber))
			check(err)
			return fullURI, indexNumber, false, pageTopicBroken
		}
	}

	defer response.Body.Close()
	if response.StatusCode != 200 {
		// the node that is linked to is broken - that is it does not exist
		indexNumber++
		indexNames[indexNumber] = shortURI

		appearanceInfo[shortURI]++
		listOfDuplicateInfo = append(listOfDuplicateInfo, duplicateInfo{
			id:        indexNumber,
			pageType:  pageTopic,
			parentURI: parentURI,
			shortURI:  shortURI,
		})
		listOfPageData = append(listOfPageData, pageData{
			id:              indexNumber,
			subSectionIndex: parentID,
			pageType:        pageTopic,
			uriStatus:       pageTopicBroken,
			shortURI:        shortURI,
			parentURI:       "https://www.ons.gov.uk" + parentURI + "/data",
			fixedPayload:    []byte{},
		})

		fmt.Printf("\nERROR on ONS website /data field: %v\n\n", response.StatusCode)
		fmt.Printf("URI does not exist:  %v\n", fullURI)
		_, err := fmt.Fprintf(graphVizFile, "    subgraph %s {\n        %s [shape = box, color=red, style=bold, label = \"%s\"]\n    }\n", gvPage, gvPage, gvPageLabel+"\n ** MISSING: /data **"+fmt.Sprintf("\\n%v", indexNumber))
		check(err)
		return fullURI, indexNumber, false, pageTopicBroken
	}
	bodyText, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		fmt.Printf("getPage: ReadAll(response.Body) failed\n")
		os.Exit(2)
	}

	pageCount++

	// The following section of code checks that the data structure 'DataResponse' has all the fields
	// needed to read in the desired web page.

	// NOTE: the process of Unmarshal and Marshal manages to 'break' the comparison i'm doing later on
	//       by (somewhere in the process) converting things like \u0027 into a single character of a single quote
	//       within some string fields.
	// SO: the following call to replaceUnicodeWithASCII() replaces things like \u0027 with whatever their one
	//     byte ASCII equivalent is (hopefully) (it seems to work for ONS web pages)
	fixedJSON := replaceUnicodeWithASCII(bodyText)

	var data DataResponse

	// Unmarshal body bytes to model
	if err := json.Unmarshal(fixedJSON, &data); err != nil {
		fmt.Println(err)
		fmt.Printf("getPage: json.Unmarshal failed\n")
		os.Exit(3)
	}

	// Marshal provided model
	payload, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		fmt.Printf("getPage: json.Marshal failed\n")
		os.Exit(4)
	}

	// This effectively checks that the struct 'DataResponse' has all the fields needed ..
	// the 'payLoad' should equal the 'fixedJSON' .. if not DataResponse needs adjusting
	if !bytes.Equal(payload, fixedJSON) {
		fmt.Printf("Processing topic page: %s\n", fullURI)
		fmt.Printf("Unmarshal / Marshal mismatch.\nInspect the saved .json files and fix struct DataResponse\n")
		_, err = fmt.Fprintf(bodyTextFile, "%s\n", fixedJSON)
		check(err)
		_, err = fmt.Fprintf(checkFile, "%s\n", payload)
		check(err)
		os.Exit(5)
	}

	var thisPageType string
	if data.Type != nil {
		thisPageType = *data.Type
		fmt.Printf("Page Type: %s\n", *data.Type) // to help understand exact page type
	} else {
		thisPageType = "Type: unknown"
		fmt.Printf("ODD (possible error in page data): no 'Type'\n")
	}

	if data.Sections == nil {
		// This is a Content page (a termination part of the graph)
		indexNumber++
		indexNames[indexNumber] = shortURI

		appearanceInfo[shortURI]++
		listOfDuplicateInfo = append(listOfDuplicateInfo, duplicateInfo{
			id:        indexNumber,
			pageType:  pageContent,
			parentURI: parentURI,
			shortURI:  shortURI,
		})
		listOfPageData = append(listOfPageData, pageData{
			id:              indexNumber,
			subSectionIndex: parentID,
			pageType:        pageContent,
			uriStatus:       pageContent,
			shortURI:        shortURI,
			parentURI:       parentURI,
			fixedPayload:    fixedJSON,
		})

		fmt.Printf("%v : Content Page: %v\n", indexNumber, fullURI)

		info := getTerminationNodeData(&data, indexNumber, fullURI)
		switch info {
		case contentUnknown:
			// no content info read, and indicate this with a Yellow border for the content termination node ..
			_, err = fmt.Fprintf(graphVizFile, "    subgraph %s {\n        %s [shape = box, color=\"#E0F020\", style=bold, label = \"%s\"]\n", gvPage, gvPage, gvPageLabel+"\n ** NO Content **"+fmt.Sprintf("\\n%v", indexNumber))
			check(err)
		case contentNone:
			// no content info exists, and indicate this with a RED background for the content termination node ..
			_, err = fmt.Fprintf(graphVizFile, "    subgraph %s {\n        %s [shape = box, color=\"#E02020\", style=filled, label = \"%s\"]\n", gvPage, gvPage, gvPageLabel+"\n ** NO Content **"+fmt.Sprintf("\\n%v", indexNumber))
			check(err)
		case contentExists:
			_, err = fmt.Fprintf(graphVizFile, "    subgraph %s {\n        %s [shape = box, color=\"#E07020\", style=filled, label = \"%s\"]\n", gvPage, gvPage, gvPageLabel+fmt.Sprintf("\\n%v - %s", indexNumber, thisPageType))
			check(err)
		}

		if info != contentUnknown {
			// write out graph content child 'broken links'
			graphContentChildBrokenLinks(gvPage, &data, graphVizFile)
		}

		// Close the subgraph:
		_, err = fmt.Fprintf(graphVizFile, "    }\n")
		check(err)

		if info != contentUnknown {
			// and ..
			// write out graph any content child 'broken links' termination node as their own subgraph's
			graphContentChildBrokenLinksSubgraph(gvPage, &data, graphVizFile)
		}

		return fullURI, indexNumber, true, pageContent
	}

	returnedIndexNumber := indexNumber // NOTE: the linter complains about this being ineffectual, but its missing the recursion through getPage()

	if len(*data.Sections) > 0 {
		indexNumber++
		indexNames[indexNumber] = shortURI

		appearanceInfo[shortURI]++
		listOfDuplicateInfo = append(listOfDuplicateInfo, duplicateInfo{
			id:        indexNumber,
			pageType:  pageTopic,
			parentURI: parentURI,
			shortURI:  shortURI,
		})
		listOfPageData = append(listOfPageData, pageData{
			id:              indexNumber,
			subSectionIndex: parentID,
			pageType:        pageTopic,
			uriStatus:       pageTopic,
			shortURI:        shortURI,
			parentURI:       parentURI,
			fixedPayload:    fixedJSON,
		})
		returnedIndexNumber = indexNumber

		fmt.Printf("%v : Topic Page: %v\n", indexNumber, fullURI)

		// do title for subgraph with children (of Topics)
		if parentURI == "" {
			// have oval shape for top level URI
			_, err = fmt.Fprintf(graphVizFile, "    subgraph %s {\n        %s [label = \"%s\"]\n\n", gvPage, gvPage, gvPageLabel+fmt.Sprintf("\\n%v", indexNumber))
			check(err)
		} else {
			// have box for sub levels (takes up less space and graph looks better)
			_, err = fmt.Fprintf(graphVizFile, "    subgraph %s {\n        %s [shape = box, label = \"%s\"]\n\n", gvPage, gvPage, gvPageLabel+fmt.Sprintf("\\n%v - %s", indexNumber, thisPageType))
			check(err)
		}

		// write out child 'sections' links
		for _, link := range *data.Sections {
			// remove leading '/'
			s := *link.URI
			gvLink := s[1:]
			// replace remaining '/' with '_' to have a connection name that is suitable for graphviz
			gvLink = strings.ReplaceAll(gvLink, "/", "_")
			_, err = fmt.Fprintf(graphVizFile, "        %s -> %s\n", gvPage, gvLink)
			check(err)
		}

		// ===========================================
		// write out child 'highlighted links'
		if data.HighlightedLinks != nil {
			if len(*data.HighlightedLinks) > 0 {
				for _, link := range *data.HighlightedLinks {
					// remove leading '/'
					s := *link.URI
					gvLink := s[1:]
					// replace remaining '/' with '_' to have a connection name that is suitable for graphviz
					gvLink = strings.ReplaceAll(gvLink, "/", "_")
					_, err = fmt.Fprintf(graphVizFile, "        %s -> %s\n", gvPage, gvLink)
					check(err)
				}
			}
		}
		// Close the subgraph:
		_, err = fmt.Fprintf(graphVizFile, "    }\n")
		check(err)

		// ===========================================
		if getNodeData(&data, returnedIndexNumber, fullURI, bodyTextFile, checkFile) {
			// This page has 'highlighted' type which has been re-maped to spotlight .. that is on a content page
			// so, Tag page as also having content:
			listOfDuplicateInfo = append(listOfDuplicateInfo, duplicateInfo{
				id:        returnedIndexNumber,
				pageType:  pageContent,
				parentURI: parentURI,
				shortURI:  shortURI,
			})
			listOfPageData = append(listOfPageData, pageData{
				id:              returnedIndexNumber,
				subSectionIndex: parentID,
				pageType:        pageContent,
				uriStatus:       pageContent,
				shortURI:        shortURI,
				parentURI:       parentURI,
				fixedPayload:    fixedJSON,
			})
		}

		// and ..
		// write out any child 'highlighted links' termination node as their own subgraph's
		if data.HighlightedLinks != nil {
			if len(*data.HighlightedLinks) > 0 {
				for _, link := range *data.HighlightedLinks {
					// remove leading '/'
					s := *link.URI
					gvLink := s[1:]
					// replace remaining '/' with '_' to have a connection name that is suitable for graphviz
					gvLink = strings.ReplaceAll(gvLink, "/", "_")
					gvLinkLabel := "/" + strings.ReplaceAll(gvLink, "_", "\\n/")
					if link.valid {
						_, err = fmt.Fprintf(graphVizFile, "    subgraph %s {\n        %s [shape = box, color=\"#30A8A0\", style=filled, label = \"%s\"]\n    }\n", gvLink, gvLink, fmt.Sprintf("Topic - Highlighted: %s\\n", link.linkType)+gvLinkLabel)
						check(err)
					} else {
						_, err = fmt.Fprintf(graphVizFile, "    subgraph %s {\n        %s [shape = box, color=\"#30A8A0\", style=bold, label = \"%s\"]\n    }\n", gvLink, gvLink, "Topic - Highlighted:\\n"+gvLinkLabel)
						check(err)
					}
				}
			}
		}

		type sectionResult struct {
			fullURI string
			id      int
			pagType allowedPageType
		}

		var sectionResults []sectionResult

		for _, link := range *data.Sections {
			// get a sub section page

			// recursively call self ..
			retFullURI, retIndex, valid, pgType := getPage(returnedIndexNumber, graphVizFile, shortURI, *link.URI)
			if valid {
				// the page is NOT broken ..
				sectionResults = append(sectionResults, sectionResult{fullURI: retFullURI, id: retIndex, pagType: pgType})
				if pgType == pageTopic {
					// save the sub topics id's
					// 'returnedIndexNumber' is the parent ID of 'retIndex'
					listOfPageData = append(listOfPageData, pageData{
						id:              retIndex,
						subSectionIndex: returnedIndexNumber,
						pageType:        pageTopicSubtopicID,
						uriStatus:       pageTopicSubtopicID,
						shortURI:        *link.URI,
						parentURI:       shortURI,
					})
				}
			}
		}
		if len(sectionResults) > 0 {
			fmt.Printf("Page: %v\n", fullURI)
			if sectionResults[0].pagType == pageTopic {
				fmt.Printf("Sections (Sub Topics)\n")
			} else {
				fmt.Printf("No Sections (Sub Content)\n")
			}
			for _, section := range sectionResults {
				fmt.Printf("  %v  :  %v\n", section.id, section.fullURI)
			}
			fmt.Printf("\n")
		}
	} else {
		// the node that is linked to is broken - no topic or content info
		indexNumber++
		indexNames[indexNumber] = shortURI

		appearanceInfo[shortURI]++
		listOfDuplicateInfo = append(listOfDuplicateInfo, duplicateInfo{
			id:        indexNumber,
			pageType:  pageTopic,
			parentURI: parentURI,
			shortURI:  shortURI,
		})
		listOfPageData = append(listOfPageData, pageData{
			id:              indexNumber,
			subSectionIndex: parentID,
			pageType:        pageTopic,
			uriStatus:       pageTopicBroken,
			shortURI:        shortURI,
			parentURI:       "https://www.ons.gov.uk" + parentURI + "/data",
			fixedPayload:    []byte{},
		})

		fmt.Printf("ERROR: page has no topic or content links ..\n%v\n", fullURI)
		_, err := fmt.Fprintf(graphVizFile, "    subgraph %s {\n        %s [shape = box, color=\"#E0F020\", style=filled, label = \"%s\"]\n    }\n", gvPage, gvPage, gvPageLabel+"\n ** MISSING: topic & content **"+fmt.Sprintf("\\n%v", indexNumber))
		check(err)
		return fullURI, indexNumber, false, pageTopicBroken
	}
	return fullURI, returnedIndexNumber, true, pageTopic
}

func contentChild(uri string, gvPage string, graphVizFile io.Writer) {
	// remove leading '/'
	gvLink := uri[1:]
	// replace remaining '/' with '_' to have a connection name that is suitable for graphviz
	gvLink = strings.ReplaceAll(gvLink, "/", "_")
	gvLink = strings.ReplaceAll(gvLink, "-", "_") // can't have minus signs in dates
	_, err := fmt.Fprintf(graphVizFile, "        %s -> %s\n", gvPage, gvLink)
	check(err)
}

func graphContentChildBrokenLinks(gvPage string, data *DataResponse, graphVizFile io.Writer) {

	// read any child 'Items                     (Timeseries) links' and save their page /data
	if data.Items != nil {
		if len(*data.Items) > 0 {
			fmt.Printf("Getting: Items\n")
			for _, link := range *data.Items {
				if !link.valid {
					contentChild(*link.URI, gvPage, graphVizFile)
				} else if cfg.GraphAllContent {
					contentChild(*link.URI, gvPage, graphVizFile)
				}
			}
		}
	}

	// read any child 'Datasets                  (Static datasets) links' and save their page /data
	if data.Datasets != nil {
		if len(*data.Datasets) > 0 {
			fmt.Printf("Getting: Datasets\n")
			for _, link := range *data.Datasets {
				if !link.valid {
					contentChild(*link.URI, gvPage, graphVizFile)
				} else if cfg.GraphAllContent {
					contentChild(*link.URI, gvPage, graphVizFile)
				}
			}
		}
	}

	// read any child 'StatsBulletins            (Bulletins) links' and save their page /data
	if data.StatsBulletins != nil {
		if len(*data.StatsBulletins) > 0 {
			fmt.Printf("Getting: StatsBulletins\n")
			for _, link := range *data.StatsBulletins {
				if !link.valid {
					contentChild(*link.URI, gvPage, graphVizFile)
				} else if cfg.GraphAllContent {
					contentChild(*link.URI, gvPage, graphVizFile)
				}
			}
		}
	}

	// read any child 'RelatedArticles           (Articles) links' and save their page /data
	if data.RelatedArticles != nil {
		if len(*data.RelatedArticles) > 0 {
			fmt.Printf("Getting: RelatedArticles\n")
			for _, link := range *data.RelatedArticles {
				if !link.valid {
					contentChild(*link.URI, gvPage, graphVizFile)
				} else if cfg.GraphAllContent {
					contentChild(*link.URI, gvPage, graphVizFile)
				}
			}
		}
	}

	// read any child 'RelatedMethodology        (Methodologies) links' and save their page /data
	if data.RelatedMethodology != nil {
		if len(*data.RelatedMethodology) > 0 {
			fmt.Printf("Getting: RelatedMethodology\n")
			for _, link := range *data.RelatedMethodology {
				if !link.valid {
					contentChild(*link.URI, gvPage, graphVizFile)
				} else if cfg.GraphAllContent {
					contentChild(*link.URI, gvPage, graphVizFile)
				}
			}
		}
	}

	// read any child 'RelatedMethodologyArticle (Methodology_articles) links' and save their page /data
	if data.RelatedMethodologyArticle != nil {
		if len(*data.RelatedMethodologyArticle) > 0 {
			fmt.Printf("Getting: RelatedMethodologyArticle\n")
			for _, link := range *data.RelatedMethodologyArticle {
				if !link.valid {
					contentChild(*link.URI, gvPage, graphVizFile)
				} else if cfg.GraphAllContent {
					contentChild(*link.URI, gvPage, graphVizFile)
				}
			}
		}
	}

	// read any child 'HighlightedContent        (Spotlight) links' and save their page /data
	if data.HighlightedContent != nil {
		if len(*data.HighlightedContent) > 0 {
			fmt.Printf("Getting: HighlightedContent\n")
			for _, link := range *data.HighlightedContent {
				if !link.valid {
					contentChild(*link.URI, gvPage, graphVizFile)
				} else if cfg.GraphAllContent {
					contentChild(*link.URI, gvPage, graphVizFile)
				}
			}
		}
	}
}

func linksBrokenSubgraph(uri string, description string, graphVizFile io.Writer) {
	// remove leading '/'
	gvLink := uri[1:]
	// replace remaining '/' with '_' to have a connection name that is suitable for graphviz
	gvLink = strings.ReplaceAll(gvLink, "/", "_")
	gvLinkLabel := "/" + strings.ReplaceAll(gvLink, "_", "\\n/")
	gvLink = strings.ReplaceAll(gvLink, "-", "_")           // can't have minus signs in dates
	gvLinkLabel = strings.ReplaceAll(gvLinkLabel, "-", "_") // can't have minus signs in dates
	_, err := fmt.Fprintf(graphVizFile, "    subgraph %s {\n        %s [shape = box, color=\"#1038E0\", style=bold, label = \"%s\"]\n    }\n", gvLink, gvLink, fmt.Sprintf("** %s - Broken URI: **\\n", description)+gvLinkLabel)
	check(err)
}

func linksGoodSubgraph(uri string, description string, graphVizFile io.Writer, linkType string) {
	// remove leading '/'
	gvLink := uri[1:]
	// replace remaining '/' with '_' to have a connection name that is suitable for graphviz
	gvLink = strings.ReplaceAll(gvLink, "/", "_")
	gvLinkLabel := "/" + strings.ReplaceAll(gvLink, "_", "\\n/")
	gvLink = strings.ReplaceAll(gvLink, "-", "_")           // can't have minus signs in dates
	gvLinkLabel = strings.ReplaceAll(gvLinkLabel, "-", "_") // can't have minus signs in dates
	var err error
	if cfg.ColourContent {
		var col string = "#E10600" // to indicate error
		// the shape of this graph is oval
		switch linkType {
		case articleCollectionName:
			col = "#F9A12E" // Radiant Yellow
		case articleDownloadCollectionName:
			col = "#FC766A" // Living Coral
		case bulletinCollectionName:
			col = "#CB4AC7" // Purple
		case compendiumDataCollectionName:
			col = "#AE0E36" // Crimson
		case compendiumLandingPageCollectionName:
			col = "#6A9AB3" // Sky Blue
		case datasetLandingPageCollectionName:
			col = "#F1AC88" // Peach
		case staticMethodologyCollectionName:
			col = "#F6EA7B" // Lemon Verbena
		case staticMethodologyDownloadCollectionName:
			col = "#E683A9" // Aurora Pink
		case staticQmiCollectionName:
			col = "#078282" // Teal
		case timeseriesCollectionName:
			col = "#00A4CC" // Out of the Blue
		case chartCollectionName:
			col = "#00A400" //
		case equationCollectionName:
			col = "#070082" //
		case imageCollectionName:
			col = "#0083A9" //
		case releaseCollectionName:
			col = "#86EA00" //
		case listCollectionName:
			col = "#81AC88" //
		case staticPageCollectionName:
			col = "#6A00B3" //
		case staticAdhocCollectionName:
			col = "#5E0E36" //
		case referenceTablesCollectionName:
			col = "#6B4A67" //
		case compendiumChapterCollectionName:
			col = "#FC006A" //
		case staticLandingPageCollectionName:
			col = "#A9712E" //
		case staticArticleCollectionName:
			col = "#F9916E" //
		case datasetCollectionName:
			col = "#39A16E" //
		case timeseriesDatasetCollectionName:
			col = "#F66A3B" //
		case taxonomyLandingPageCollectionName:
			col = "#D6CA5B" //
		default:
			linkType = "ERROR - unknown: " + linkType
		}
		_, err = fmt.Fprintf(graphVizFile, "    subgraph %s {\n        %s [color=\"%s\", style=filled, label = \"%s\"]\n    }\n", gvLink, gvLink, col, linkType+"\\n"+fmt.Sprintf("** %s - Good URI: **\\n", description)+gvLinkLabel)
	} else {
		_, err = fmt.Fprintf(graphVizFile, "    subgraph %s {\n        %s [color=\"#1038E0\", style=filled, label = \"%s\"]\n    }\n", gvLink, gvLink, fmt.Sprintf("** %s - Good URI: **\\n", description)+gvLinkLabel)
	}
	check(err)
}

func graphContentChildBrokenLinksSubgraph(gvPage string, data *DataResponse, graphVizFile io.Writer) {

	// read any child 'Items                     (Timeseries) links' and save their page /data
	if data.Items != nil {
		if len(*data.Items) > 0 {
			fmt.Printf("Getting: Items\n")
			for _, link := range *data.Items {
				if !link.valid {
					linksBrokenSubgraph(*link.URI, "Items", graphVizFile)
				} else if cfg.GraphAllContent {
					linksGoodSubgraph(*link.URI, "Items", graphVizFile, link.linkType)
				}
			}
		}
	}

	// read any child 'Datasets                  (Static datasets) links' and save their page /data
	if data.Datasets != nil {
		if len(*data.Datasets) > 0 {
			fmt.Printf("Getting: Datasets\n")
			for _, link := range *data.Datasets {
				if !link.valid {
					linksBrokenSubgraph(*link.URI, "Datasets", graphVizFile)
				} else if cfg.GraphAllContent {
					linksGoodSubgraph(*link.URI, "Datasets", graphVizFile, link.linkType)
				}
			}
		}
	}

	// read any child 'StatsBulletins            (Bulletins) links' and save their page /data
	if data.StatsBulletins != nil {
		if len(*data.StatsBulletins) > 0 {
			fmt.Printf("Getting: StatsBulletins\n")
			for _, link := range *data.StatsBulletins {
				if !link.valid {
					linksBrokenSubgraph(*link.URI, "StatsBulletins", graphVizFile)
				} else if cfg.GraphAllContent {
					linksGoodSubgraph(*link.URI, "StatsBulletins", graphVizFile, link.linkType)
				}
			}
		}
	}

	// read any child 'RelatedArticles           (Articles) links' and save their page /data
	if data.RelatedArticles != nil {
		if len(*data.RelatedArticles) > 0 {
			fmt.Printf("Getting: RelatedArticles\n")
			for _, link := range *data.RelatedArticles {
				if !link.valid {
					linksBrokenSubgraph(*link.URI, "RelatedArticles", graphVizFile)
				} else if cfg.GraphAllContent {
					linksGoodSubgraph(*link.URI, "RelatedArticles", graphVizFile, link.linkType)
				}
			}
		}
	}

	// read any child 'RelatedMethodology        (Methodologies) links' and save their page /data
	if data.RelatedMethodology != nil {
		if len(*data.RelatedMethodology) > 0 {
			fmt.Printf("Getting: RelatedMethodology\n")
			for _, link := range *data.RelatedMethodology {
				if !link.valid {
					linksBrokenSubgraph(*link.URI, "RelatedMethodology", graphVizFile)
				} else if cfg.GraphAllContent {
					linksGoodSubgraph(*link.URI, "RelatedMethodology", graphVizFile, link.linkType)
				}
			}
		}
	}

	// read any child 'RelatedMethodologyArticle (Methodology_articles) links' and save their page /data
	if data.RelatedMethodologyArticle != nil {
		if len(*data.RelatedMethodologyArticle) > 0 {
			fmt.Printf("Getting: RelatedMethodologyArticle\n")
			for _, link := range *data.RelatedMethodologyArticle {
				if !link.valid {
					linksBrokenSubgraph(*link.URI, "RelatedMethodologyArticle", graphVizFile)
				} else if cfg.GraphAllContent {
					linksGoodSubgraph(*link.URI, "RelatedMethodologyArticle", graphVizFile, link.linkType)
				}
			}
		}
	}

	// read any child 'HighlightedContent        (Spotlight) links' and save their page /data
	if data.HighlightedContent != nil {
		if len(*data.HighlightedContent) > 0 {
			fmt.Printf("Getting: HighlightedContent\n")
			for _, link := range *data.HighlightedContent {
				if !link.valid {
					linksBrokenSubgraph(*link.URI, "HighlightedContent", graphVizFile)
				} else if cfg.GraphAllContent {
					linksGoodSubgraph(*link.URI, "HighlightedContent", graphVizFile, link.linkType)
				}
			}
		}
	}
}

var articleCount int
var articleDownloadCount int
var bulletinCount int
var compendiumDataCount int
var compendiumLandingPageCount int
var datasetLandingPageCount int
var staticMethodologyCount int
var staticMethodologyDownloadCount int
var staticQmiCount int
var timeseriesCount int
var chartCount int
var productPageCount int
var tableCount int
var equationCount int
var imageCount int
var releaseCount int
var listCount int
var staticPageCount int
var staticAdhocCount int
var referenceTablesCount int
var compendiumChapterCount int
var staticLandingPageCount int
var staticArticleCount int
var datasetCount int
var timeseriesDatasetCount int
var taxonomyLandingPageCount int

// store the shortUri and count to prevent processing a page more than once
var contentDuplicateCheck = make(map[string]int) // key: shortURI, value: 1 or more indicates already saved
var uriCollectionName = make(map[string]string)  // key: shortURI, value: the name of the colection that the URI is storred in

func saveContentPageToCollection(collectionJsFile *os.File, id *int, collectionName string, bodyTextCopy []byte, shortURI string) {
	// The original /data endpoint information read has passed tests so now we
	// write out the original json code together with an extra 'id'

	// NOTE: splitting the content info out to separate collections as per the content 'Type' may not
	//       be whats needed, but its a good start for demonstrating that the content has been
	//       extracted without error.

	if contentDuplicateCheck[shortURI] > 0 {
		// keep incrementing the duplication count (in case its of use at some point)
		contentDuplicateCheck[shortURI]++
		return
	}
	contentDuplicateCheck[shortURI]++
	uriCollectionName[shortURI] = collectionName

	*id++
	idStr := strconv.Itoa(*id)

	_, err := fmt.Fprintf(collectionJsFile, "db."+collectionName+".insertOne({")
	check(err)

	// write out an 'id' for this data file
	_, err = fmt.Fprintf(collectionJsFile, "\n    \"id\": \""+idStr+"\",\n")
	check(err)

	// write out what should be a unique key that can be indexed on ..
	_, err = fmt.Fprintf(collectionJsFile, "    \"id_uri\": \""+shortURI+"\",\n")
	check(err)

	// Strip out the first character which is an opening curly brace so that we get a correctly formed
	// java script line
	_, err = collectionJsFile.Write(bodyTextCopy[1:])
	check(err)

	_, err = fmt.Fprintf(collectionJsFile, ")\n")
	check(err)
}

func addSections(uriList *[]string, field *[]sections) {
	if field != nil {
		if len(*field) > 0 {
			for _, info := range *field {
				if info.URI != nil {
					*uriList = append(*uriList, *info.URI)
				}
			}
		}
	}
}

func addRelatedData(uriList *[]string, field *[]relatedData) {
	if field != nil {
		if len(*field) > 0 {
			for _, info := range *field {
				if info.URI != nil {
					*uriList = append(*uriList, *info.URI)
				}
			}
		}
	}
}

func addRelatedDocuments(uriList *[]string, field *[]relatedDocuments) {
	if field != nil {
		if len(*field) > 0 {
			for _, info := range *field {
				if info.URI != nil {
					*uriList = append(*uriList, *info.URI)
				}
			}
		}
	}
}

func addCharts(uriList *[]string, field *[]charts) {
	if field != nil {
		if len(*field) > 0 {
			for _, info := range *field {
				if info.URI != nil {
					*uriList = append(*uriList, *info.URI)
				}
			}
		}
	}
}

func addTables(uriList *[]string, field *[]tables) {
	if field != nil {
		if len(*field) > 0 {
			for _, info := range *field {
				if info.URI != nil {
					*uriList = append(*uriList, *info.URI)
				}
			}
		}
	}
}

func addImages(uriList *[]string, field *[]images) {
	// We can't read an image, so we don't check that the link is OK ..
	// and thus this code is commented out
	// If a way could be found to check that a link to a .png or .jpg is OK
	// then this could be put back in:
	/*	if field != nil {
		if len(*field) > 0 {
			for _, info := range *field {
				if info.URI != nil {
					*uriList = append(*uriList, *info.URI)
				}
			}
		}
	}*/
}

func addEquations(uriList *[]string, field *[]equations) {
	if field != nil {
		if len(*field) > 0 {
			for _, info := range *field {
				if info.URI != nil {
					*uriList = append(*uriList, *info.URI)
				}
			}
		}
	}
}

func addLinks(uriList *[]string, field *[]links) {
	if field != nil {
		if len(*field) > 0 {
			for _, info := range *field {
				if info.URI != nil {
					*uriList = append(*uriList, *info.URI)
				}
			}
		}
	}
}

func addRelatedMethodology(uriList *[]string, field *[]relatedMethodology) {
	if field != nil {
		if len(*field) > 0 {
			for _, info := range *field {
				if info.URI != nil {
					*uriList = append(*uriList, *info.URI)
				}
			}
		}
	}
}

func addRelatedMethodologyArticle(uriList *[]string, field *[]relatedMethodologyArticle) {
	if field != nil {
		if len(*field) > 0 {
			for _, info := range *field {
				if info.URI != nil {
					*uriList = append(*uriList, *info.URI)
				}
			}
		}
	}
}

func addVersions(uriList *[]string, field *[]versions) {
	if field != nil {
		if len(*field) > 0 {
			for _, info := range *field {
				if info.URI != nil {
					*uriList = append(*uriList, *info.URI)
				}
			}
		}
	}
}

func addTopics(uriList *[]string, field *[]ctopics) {
	if field != nil {
		if len(*field) > 0 {
			for _, info := range *field {
				if info.URI != nil {
					*uriList = append(*uriList, *info.URI)
				}
			}
		}
	}
}

func addRelatedDatasets(uriList *[]string, field *[]relatedDatasets) {
	if field != nil {
		if len(*field) > 0 {
			for _, info := range *field {
				if info.URI != nil {
					*uriList = append(*uriList, *info.URI)
				}
			}
		}
	}
}

func addDatasets(uriList *[]string, field *[]datasets) {
	if field != nil {
		if len(*field) > 0 {
			for _, info := range *field {
				if info.URI != nil {
					*uriList = append(*uriList, *info.URI)
				}
			}
		}
	}
}

func addChapters(uriList *[]string, field *[]chapters) {
	if field != nil {
		if len(*field) > 0 {
			for _, info := range *field {
				if info.URI != nil {
					*uriList = append(*uriList, *info.URI)
				}
			}
		}
	}
}

func addRelatedFilterableDatasets(uriList *[]string, field *[]relatedFilterableDatasets) {
	if field != nil {
		if len(*field) > 0 {
			for _, info := range *field {
				if info.URI != nil {
					if !strings.HasPrefix(*info.URI, "/datasets/") {
						// This is NOT a 'Choose My Data' page, so add it
						// (we skip CMD pages because they do not have a '/data' suffix indicating
						//  this is a page that can not be processed)
						*uriList = append(*uriList, *info.URI)
					}
				}
			}
		}
	}
}

func addSourceDatasets(uriList *[]string, field *[]sourceDatasets) {
	if field != nil {
		if len(*field) > 0 {
			for _, info := range *field {
				if info.URI != nil {
					*uriList = append(*uriList, *info.URI)
				}
			}
		}
	}
}

func getURIListFromArticle(containintURI string, data *articleResponse) []string {
	var uriList []string

	addSections(&uriList, data.Sections)
	addRelatedData(&uriList, data.RelatedData)
	addRelatedDocuments(&uriList, data.RelatedDocuments)
	addCharts(&uriList, data.Charts)
	addTables(&uriList, data.Tables)
	addImages(&uriList, data.Images)
	addEquations(&uriList, data.Equations)
	addLinks(&uriList, data.Links)
	addRelatedMethodology(&uriList, data.RelatedMethodology)
	addRelatedMethodologyArticle(&uriList, data.RelatedMethodologyArticle)
	addVersions(&uriList, data.Versions)
	addTopics(&uriList, data.Topics)

	return uriList
}

func getURIListFromArticleDownload(containintURI string, data *articleDownloadResponse) []string {
	var uriList []string

	if cfg.OnlyFirstFullDepth {
		return uriList
	}

	addRelatedData(&uriList, data.RelatedData)
	addRelatedDocuments(&uriList, data.RelatedDocuments)
	addCharts(&uriList, data.Charts)
	addTables(&uriList, data.Tables)
	addImages(&uriList, data.Images)
	addEquations(&uriList, data.Equations)
	addLinks(&uriList, data.Links)
	addRelatedMethodology(&uriList, data.RelatedMethodology)
	addRelatedMethodologyArticle(&uriList, data.RelatedMethodologyArticle)
	addVersions(&uriList, data.Versions)
	addTopics(&uriList, data.Topics)

	return uriList
}

func getURIListFromBulletin(containintURI string, data *bulletinResponse) []string {
	var uriList []string

	if cfg.OnlyFirstFullDepth {
		return uriList
	}

	addSections(&uriList, data.Sections)
	addRelatedData(&uriList, data.RelatedData)
	addRelatedDocuments(&uriList, data.RelatedDocuments)
	addCharts(&uriList, data.Charts)
	addTables(&uriList, data.Tables)
	addImages(&uriList, data.Images)
	addEquations(&uriList, data.Equations)
	addLinks(&uriList, data.Links)
	addRelatedMethodology(&uriList, data.RelatedMethodology)
	addRelatedMethodologyArticle(&uriList, data.RelatedMethodologyArticle)
	addVersions(&uriList, data.Versions)
	addTopics(&uriList, data.Topics)

	return uriList
}

func getURIListFromCompendiumData(containintURI string, data *compendiumDataResponse) []string {
	var uriList []string

	if cfg.OnlyFirstFullDepth {
		return uriList
	}

	addRelatedDatasets(&uriList, data.RelatedDatasets)
	addRelatedDocuments(&uriList, data.RelatedDocuments)
	addRelatedMethodology(&uriList, data.RelatedMethodology)
	addRelatedMethodologyArticle(&uriList, data.RelatedMethodologyArticle)
	addVersions(&uriList, data.Versions)
	addTopics(&uriList, data.Topics)

	return uriList
}

func getURIListFromCompendiumLandingPage(containintURI string, data *compendiumLandingPageResponse) []string {
	var uriList []string

	if cfg.OnlyFirstFullDepth {
		return uriList
	}

	addDatasets(&uriList, data.Datasets)
	addChapters(&uriList, data.Chapters)
	addRelatedDocuments(&uriList, data.RelatedDocuments)
	addRelatedData(&uriList, data.RelatedData)
	addRelatedMethodology(&uriList, data.RelatedMethodology)
	addRelatedMethodologyArticle(&uriList, data.RelatedMethodologyArticle)
	addTopics(&uriList, data.Topics)

	return uriList
}

func getURIListFromDatasetLandingPage(containintURI string, data *datasetLandingPageResponse) []string {
	var uriList []string

	if cfg.OnlyFirstFullDepth {
		return uriList
	}

	addRelatedFilterableDatasets(&uriList, data.RelatedFilterableDatasets)
	addRelatedDatasets(&uriList, data.RelatedDatasets)
	addRelatedDocuments(&uriList, data.RelatedDocuments)
	addDatasets(&uriList, data.Datasets)
	addLinks(&uriList, data.Links)
	addRelatedMethodology(&uriList, data.RelatedMethodology)
	addRelatedMethodologyArticle(&uriList, data.RelatedMethodologyArticle)
	addVersions(&uriList, data.Versions)
	addTopics(&uriList, data.Topics)

	return uriList
}

func getURIListFromStaticMethodology(containintURI string, data *staticMethodologyResponse) []string {
	var uriList []string

	if cfg.OnlyFirstFullDepth {
		return uriList
	}

	addRelatedDocuments(&uriList, data.RelatedDocuments)
	addLinks(&uriList, data.Links)
	addSections(&uriList, data.Sections)
	addRelatedData(&uriList, data.RelatedData)
	addCharts(&uriList, data.Charts)
	addTables(&uriList, data.Tables)
	addImages(&uriList, data.Images)
	addEquations(&uriList, data.Equations)
	addRelatedMethodology(&uriList, data.RelatedMethodology)
	addRelatedMethodologyArticle(&uriList, data.RelatedMethodologyArticle)
	addVersions(&uriList, data.Versions)
	addTopics(&uriList, data.Topics)

	return uriList
}

func getURIListFromStaticMethodologyDownload(containintURI string, data *staticMethodologyDownloadResponse) []string {
	var uriList []string

	if cfg.OnlyFirstFullDepth {
		return uriList
	}

	addRelatedDocuments(&uriList, data.RelatedDocuments)
	addRelatedDatasets(&uriList, data.RelatedDatasets)
	addLinks(&uriList, data.Links)
	addTopics(&uriList, data.Topics)

	return uriList
}

func getURIListFromStaticQmi(containintURI string, data *staticQmiResponse) []string {
	var uriList []string

	if cfg.OnlyFirstFullDepth {
		return uriList
	}

	addRelatedDocuments(&uriList, data.RelatedDocuments)
	addRelatedDatasets(&uriList, data.RelatedDatasets)
	addLinks(&uriList, data.Links)
	addTopics(&uriList, data.Topics)

	return uriList
}

func getURIListFromTimeseries(containintURI string, data *timeseriesResponse) []string {
	var uriList []string

	if cfg.OnlyFirstFullDepth {
		return uriList
	}

	addSourceDatasets(&uriList, data.SourceDatasets)
	addRelatedDatasets(&uriList, data.RelatedDatasets)
	addRelatedDocuments(&uriList, data.RelatedDocuments)
	addRelatedData(&uriList, data.RelatedData)
	addVersions(&uriList, data.Versions)
	addTopics(&uriList, data.Topics)

	return uriList
}

func getURIListFromRelease(containintURI string, data *releaseResponse) []string {
	var uriList []string

	if cfg.OnlyFirstFullDepth {
		return uriList
	}

	addRelatedDocuments(&uriList, data.RelatedDocuments)
	addRelatedDatasets(&uriList, data.RelatedDatasets)
	addRelatedMethodology(&uriList, data.RelatedMethodology)
	addRelatedMethodologyArticle(&uriList, data.RelatedMethodologyArticle)
	addLinks(&uriList, data.Links)

	return uriList
}

func getURIListFromStaticPage(containintURI string, data *staticPageResponse) []string {
	var uriList []string

	if cfg.OnlyFirstFullDepth {
		return uriList
	}

	addLinks(&uriList, data.Links)

	return uriList
}

func getURIListFromStaticAdhoc(containintURI string, data *staticAdhocResponse) []string {
	var uriList []string

	if cfg.OnlyFirstFullDepth {
		return uriList
	}

	addLinks(&uriList, data.Links)

	return uriList
}

func getURIListFromReferenceTables(containintURI string, data *referenceTablesResponse) []string {
	var uriList []string

	if cfg.OnlyFirstFullDepth {
		return uriList
	}

	addRelatedDocuments(&uriList, data.RelatedDocuments)
	addRelatedMethodology(&uriList, data.RelatedMethodology)

	return uriList
}

func getURIListFromCompendiumChapter(containintURI string, data *compendiumChapterResponse) []string {
	var uriList []string

	if cfg.OnlyFirstFullDepth {
		return uriList
	}

	addSections(&uriList, data.Sections)
	addRelatedData(&uriList, data.RelatedData)
	addRelatedDocuments(&uriList, data.RelatedDocuments)
	addCharts(&uriList, data.Charts)
	addTables(&uriList, data.Tables)
	addImages(&uriList, data.Images)
	addEquations(&uriList, data.Equations)
	addLinks(&uriList, data.Links)
	addRelatedMethodology(&uriList, data.RelatedMethodology)
	addRelatedMethodologyArticle(&uriList, data.RelatedMethodologyArticle)
	addVersions(&uriList, data.Versions)

	return uriList
}

func getURIListFromStaticLandingPage(containintURI string, data *staticLandingPageResponse) []string {
	var uriList []string

	if cfg.OnlyFirstFullDepth {
		return uriList
	}

	addSections(&uriList, data.Sections)
	addLinks(&uriList, data.Links)

	return uriList
}

func getURIListFromStaticArticle(containintURI string, data *staticArticleResponse) []string {
	var uriList []string

	if cfg.OnlyFirstFullDepth {
		return uriList
	}

	addLinks(&uriList, data.Links)
	addSections(&uriList, data.Sections)
	addCharts(&uriList, data.Charts)
	addTables(&uriList, data.Tables)
	addImages(&uriList, data.Images)
	addEquations(&uriList, data.Equations)

	return uriList
}

func getURIListFromDataset(containintURI string, data *datasetResponse) []string {
	var uriList []string

	if cfg.OnlyFirstFullDepth {
		return uriList
	}

	addVersions(&uriList, data.Versions)

	return uriList
}

func getURIListFromTimeseriesDataset(containintURI string, data *timeseriesDatasetResponse) []string {
	var uriList []string

	if cfg.OnlyFirstFullDepth {
		return uriList
	}

	addVersions(&uriList, data.Versions)

	return uriList
}

var depth int = 1
var maxDepth = depth

func unmarshalFail(uri string, err error, location int) {
	fmt.Printf("fullURI: %s\n", uri)
	fmt.Println(err)
	fmt.Printf("getPageData: json.Unmarshal failed %v\n", location)
	os.Exit(100)
}

func marshalFail(uri string, err error, location int) {
	fmt.Printf("fullURI: %s\n", uri)
	fmt.Println(err)
	fmt.Printf("getPageData: json.Marshal failed %v\n", location)
	os.Exit(101)
}

func checkMarshaling(fullURI string, err error, location int, payload *[]byte, fixedJSON *[]byte, structName string) {
	if err != nil {
		marshalFail(fullURI, err, location)
	}
	fixedPayloadJSON := replaceUnicodeWithASCII(*payload)

	// This effectively checks that the struct 'structName' has all the fields needed ..
	// the 'payLoad' should equal the 'fixedJSON' .. if not structName needs adjusting
	if !bytes.Equal(fixedPayloadJSON, *fixedJSON) {
		fmt.Printf("Processing page: %s\n", fullURI)
		fmt.Printf("Unmarshal / Marshal mismatch - %v.\nInspect the saved .json files and fix struct %s\n", location, structName)
		_, err = fmt.Fprintf(bodyTextFile, "%s\n", fixedJSON)
		check(err)
		_, err = fmt.Fprintf(checkFile, "%s\n", fixedPayloadJSON)
		check(err)
		os.Exit(102)
	}
}

func checkMarshalingDeepEqual(fullURI string, err error, location int, payload *[]byte, fixedJSON *[]byte, structName string) {
	if err != nil {
		marshalFail(fullURI, err, location)
	}
	fixedPayloadJSON := replaceUnicodeWithASCII(*payload)

	// This effectively checks that the struct 'structName' has all the fields needed ..
	// the 'payLoad' should equal the 'fixedJSON' .. if not structName needs adjusting
	if !bytes.Equal(fixedPayloadJSON, *fixedJSON) {
		// The binary comparison will typically fail for struct 'chartResponse'
		// because it contains map[string]string which after unmarshaling and marshaling ..
		// items in the maps may not in the same order.

		// So, we do a an unraveling of the binary JSON to lines of text, sort and then compare ..

		var prettyJSON1 bytes.Buffer
		err = json.Indent(&prettyJSON1, fixedPayloadJSON, "", "    ")
		check(err) // should nt get an error, but just in case

		var prettyJSON2 bytes.Buffer
		err = json.Indent(&prettyJSON2, *fixedJSON, "", "    ")
		check(err) // should not get an error, but just in case

		line1 := strings.Split(prettyJSON1.String(), "\n")
		line2 := strings.Split(prettyJSON2.String(), "\n")

		sort.Strings(line1)
		sort.Strings(line2)

		// maps don't have their fields sorted which results in otherwise equal lines having and not having commas
		// on the end of them, so to allow the DeepEqual below to work, the commas on the ends of the lines need removing.
		for i := 0; i < len(line1); i++ {
			line1[i] = strings.TrimSuffix(line1[i], ",")
		}
		for i := 0; i < len(line2); i++ {
			line2[i] = strings.TrimSuffix(line2[i], ",")
		}

		if !reflect.DeepEqual(line1, line2) {
			fmt.Printf("DeepEqual comparison failed\n")
			fmt.Printf("Processing page: %s\n", fullURI)
			fmt.Printf("Unmarshal / Marshal mismatch - %v.\nInspect the saved .json files and fix struct %s\n", location, structName)
			fmt.Printf("It helps to open these files in vscode and right click in file, select format Docuemnt\n")
			fmt.Printf(" and then save each document and then do a file comparison in an App like meld.")
			// NOTE: In the files, ignore the '&' character at the begining of one of the lines as this is just
			//       a result of one of the unmarshaled JSON lines being a pointer.
			_, err = fmt.Fprintf(checkFile, "%s\n", fixedPayloadJSON)
			check(err)
			_, err = fmt.Fprintf(checkFile, "%s\n", line1)
			check(err)
			_, err = fmt.Fprintf(bodyTextFile, "%s\n", fixedJSON)
			check(err)
			_, err = fmt.Fprintf(bodyTextFile, "%s\n", line2)
			check(err)
			os.Exit(103)
		}
	}
}

const noTitle string = "** no title **"

func extractTitle(title *string) string {
	if title != nil {
		return *title
	}
	return noTitle
}

const noDescription string = "** no description **"

func extractDescription(description *string) string {
	if description != nil {
		return *description
	}
	return noDescription
}

func getPageData(shortURI string, parentTopicNumber int, pType allowedPageType, parentURI string, index int) (int, string) {
	if cfg.FullDepth {
		if contentDuplicateCheck[shortURI] > 0 {
			// strange we've seen this link before and filtering elsewhere did not catch it.
			fmt.Printf("Repeat link: %s\n", shortURI)
			return 503, ""
		}
	}

	// Add prefix and '/data' to shortURI name
	//	fullURI := "https://www.ons.gov.uk" + shortURI + "/data"
	fullURI := "https://www.production.onsdigital.co.uk" + shortURI + "/data"

	attemptedGetCount++
	if cfg.PlayNice {
		// a little delay to play nice with ONS site and 'hopefully' not have cloudflare 'reset' the connection
		time.Sleep(100 * time.Millisecond)
	}
	response, err := http.Get(fullURI)
	if err != nil {
		fmt.Println(err)
		fmt.Printf("getPageData: http.Get(fullURI) failed\n")
		fmt.Printf("We now fabricate the response code to a 429 to instigate a retry after a delay 2\n")
		return 429, ""
	}

	defer response.Body.Close()
	if response.StatusCode != 200 {
		if response.StatusCode != 404 {
			if response.StatusCode != 429 {
				// a 503 is being seen at this point .. (it could be some other error, but whatever it is we do error action)
				fmt.Printf("\nERROR on ONS website /data field: %v\n\n", response.StatusCode)
				fmt.Printf("URI does not exist:  %v\n", fullURI)

				appearanceInfo[shortURI]++
				listOfDuplicateInfo = append(listOfDuplicateInfo, duplicateInfo{
					id:        parentTopicNumber,
					pageType:  pType,
					parentURI: parentURI,
					shortURI:  shortURI,
				})
				listOfPageData = append(listOfPageData, pageData{
					id:              parentTopicNumber,
					subSectionIndex: index,
					pageType:        pType,
					uriStatus:       pageBroken,
					shortURI:        shortURI,
					parentURI:       parentURI,
					fixedPayload:    []byte{},
				})
			} else {
				fmt.Printf("\nToo many requests\n")
				// caller will call this function again for a 429
			}
		}
		return response.StatusCode, ""
	}
	bodyText, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		fmt.Printf("getPageData: RealAll failed\n")
		os.Exit(7)
	}

	// Take a copy into another block of memory before the call to replaceUnicodeWithASCII()
	// strips out the unicode characters.
	// .. thus retaining any unicode to write back out after checks made.
	var bodyTextCopy []byte = make([]byte, len(bodyText))
	copy(bodyTextCopy, bodyText)

	fixedJSON := replaceUnicodeWithASCII(bodyText)

	// Create a list of URIs
	var uriList []string

	var title, description string

	var shape pageShape
	// Unmarshal body bytes to model
	if err := json.Unmarshal(fixedJSON, &shape); err != nil {
		fmt.Println(err)
		fmt.Printf("getPageData: json.Unmarshal failed 1\n")
		// we can get here from:
		// /employmentandlabourmarket/peopleinwork/workplacepensions#publications
		// where the '#publications' jumps one some way into a page that has already been processed
		// (though #'s are now caught elsewhere)
		//
		// OR
		//    from:
		//    https://www.ons.gov.uk/search?q=interactivetool
		//    which is a valid page, but has no data structure to extract info from
		if shortURI == "/search?q=interactivetool" {
			// NOTE: the number of exception pages may need to grow
			// say the page is unavailable
			return 503, ""
		}
		fmt.Printf("Unknown problem on page: %s\n", fullURI)
		fmt.Printf("shortURI: %s\n", shortURI)
		os.Exit(8)
	}

	// Decode each content page into a specific structure according to the 'Type' of the page ..
	// NOTE: This is done to ensure that the structure definitions are fully defined to read ALL
	//       the info in the /data endpoint.
	collectionName := *shape.Type
	switch collectionName {
	case "article":
		// "article" is linked to from these pageType on topics or content nodes
		// nodeHighlightedLinks
		// contentRelatedArticles
		// contentHighlightedContent
		var data articleResponse

		// sanity check the page has some fields for later use

		// Unmarshal body bytes to model
		if err := json.Unmarshal(fixedJSON, &data); err != nil {
			unmarshalFail(fullURI, err, 2)
		}

		// Marshal provided model
		payload, err := json.Marshal(data)
		checkMarshaling(fullURI, err, 2, &payload, &fixedJSON, "article")

		title = extractTitle(data.Description.Title)
		description = extractDescription(data.Description.MetaDescription)

		saveContentPageToCollection(articleJsFile, &articleCount, collectionName, bodyTextCopy, shortURI)
		if cfg.FullDepth {
			uriList = getURIListFromArticle(fullURI, &data)
		}

	case articleDownloadCollectionName:
		// articleDownloadCollectionName is linked to from these pageType on topics or content nodes
		// contentRelatedArticles
		var data articleDownloadResponse

		if err := json.Unmarshal(fixedJSON, &data); err != nil {
			unmarshalFail(fullURI, err, 3)
		}
		payload, err := json.Marshal(data)
		checkMarshaling(fullURI, err, 3, &payload, &fixedJSON, articleDownloadCollectionName)

		title = extractTitle(data.Description.Title)
		description = extractDescription(data.Description.MetaDescription)

		saveContentPageToCollection(articleDownloadJsFile, &articleDownloadCount, collectionName, bodyTextCopy, shortURI)
		if cfg.FullDepth {
			uriList = getURIListFromArticleDownload(fullURI, &data)
		}

	case bulletinCollectionName:
		// bulletinCollectionName is linked to from these pageType on topics or content nodes
		// nodeHighlightedLinks
		// contentStatsBulletins
		// contentHighlightedContent
		var data bulletinResponse

		if err := json.Unmarshal(fixedJSON, &data); err != nil {
			unmarshalFail(fullURI, err, 4)
		}
		payload, err := json.Marshal(data)
		checkMarshaling(fullURI, err, 4, &payload, &fixedJSON, bulletinCollectionName)

		title = extractTitle(data.Description.Title)
		description = extractDescription(data.Description.MetaDescription)

		saveContentPageToCollection(bulletinJsFile, &bulletinCount, collectionName, bodyTextCopy, shortURI)
		if cfg.FullDepth {
			uriList = getURIListFromBulletin(fullURI, &data)
		}

	case compendiumDataCollectionName:
		// compendiumDataCollectionName is linked to from these pageType on topics or content nodes
		// contentDatasets
		var data compendiumDataResponse

		if err := json.Unmarshal(fixedJSON, &data); err != nil {
			unmarshalFail(fullURI, err, 5)
		}
		payload, err := json.Marshal(data)
		checkMarshaling(fullURI, err, 5, &payload, &fixedJSON, compendiumDataCollectionName)

		title = extractTitle(data.Description.Title)
		description = extractDescription(data.Description.MetaDescription)

		saveContentPageToCollection(compendiumDataJsFile, &compendiumDataCount, collectionName, bodyTextCopy, shortURI)
		if cfg.FullDepth {
			uriList = getURIListFromCompendiumData(fullURI, &data)
		}

	case compendiumLandingPageCollectionName:
		// compendiumLandingPageCollectionName is linked to from these pageType on topics or content nodes
		// nodeHighlightedLinks
		// contentRelatedArticles
		var data compendiumLandingPageResponse

		// sanity check the page has some fields for later use

		// Unmarshal body bytes to model
		if err := json.Unmarshal(fixedJSON, &data); err != nil {
			unmarshalFail(fullURI, err, 6)
		}

		// Marshal provided model
		payload, err := json.Marshal(data)
		checkMarshaling(fullURI, err, 6, &payload, &fixedJSON, compendiumLandingPageCollectionName)

		title = extractTitle(data.Description.Title)
		description = extractDescription(data.Description.MetaDescription)

		saveContentPageToCollection(compendiumLandingPageJsFile, &compendiumLandingPageCount, collectionName, bodyTextCopy, shortURI)
		if cfg.FullDepth {
			uriList = getURIListFromCompendiumLandingPage(fullURI, &data)
		}

	case datasetLandingPageCollectionName:
		// datasetLandingPageCollectionName is linked to from these pageType on topics or content nodes
		// contentDatasets
		var data datasetLandingPageResponse

		if err := json.Unmarshal(fixedJSON, &data); err != nil {
			unmarshalFail(fullURI, err, 7)
		}
		payload, err := json.Marshal(data)
		checkMarshaling(fullURI, err, 7, &payload, &fixedJSON, datasetLandingPageCollectionName)

		title = extractTitle(data.Description.Title)
		description = extractDescription(data.Description.MetaDescription)

		saveContentPageToCollection(datasetLandingPageJsFile, &datasetLandingPageCount, collectionName, bodyTextCopy, shortURI)
		if cfg.FullDepth {
			uriList = getURIListFromDatasetLandingPage(fullURI, &data)
		}

	case staticMethodologyCollectionName:
		// staticMethodologyCollectionName is linked to from these pageType on topics or content nodes
		// contentRelatedMethodology
		// contentRelatedMethodologyArticle
		var data staticMethodologyResponse

		if err := json.Unmarshal(fixedJSON, &data); err != nil {
			unmarshalFail(fullURI, err, 8)
		}
		payload, err := json.Marshal(data)
		checkMarshaling(fullURI, err, 8, &payload, &fixedJSON, staticMethodologyCollectionName)

		title = extractTitle(data.Description.Title)
		description = extractDescription(data.Description.MetaDescription)

		saveContentPageToCollection(staticMethodologyJsFile, &staticMethodologyCount, collectionName, bodyTextCopy, shortURI)
		if cfg.FullDepth {
			uriList = getURIListFromStaticMethodology(fullURI, &data)
		}

	case staticMethodologyDownloadCollectionName:
		// staticMethodologyDownloadCollectionName is linked to from these pageType on topics or content nodes
		// contentRelatedMethodologyArticle
		var data staticMethodologyDownloadResponse

		if err := json.Unmarshal(fixedJSON, &data); err != nil {
			unmarshalFail(fullURI, err, 9)
		}
		payload, err := json.Marshal(data)
		checkMarshaling(fullURI, err, 9, &payload, &fixedJSON, staticMethodologyDownloadCollectionName)

		title = extractTitle(data.Description.Title)
		description = extractDescription(data.Description.MetaDescription)

		saveContentPageToCollection(staticMethodologyDownloadJsFile, &staticMethodologyDownloadCount, collectionName, bodyTextCopy, shortURI)
		if cfg.FullDepth {
			uriList = getURIListFromStaticMethodologyDownload(fullURI, &data)
		}

	case staticQmiCollectionName:
		// staticQmiCollectionName is linked to from these pageType on topics or content nodes
		// contentRelatedMethodology
		var data staticQmiResponse

		if err := json.Unmarshal(fixedJSON, &data); err != nil {
			unmarshalFail(fullURI, err, 10)
		}
		payload, err := json.Marshal(data)
		checkMarshaling(fullURI, err, 10, &payload, &fixedJSON, staticQmiCollectionName)

		title = extractTitle(data.Description.Title)
		description = extractDescription(data.Description.MetaDescription)

		saveContentPageToCollection(staticQmiJsFile, &staticQmiCount, collectionName, bodyTextCopy, shortURI)
		if cfg.FullDepth {
			uriList = getURIListFromStaticQmi(fullURI, &data)
		}

	case timeseriesCollectionName:
		// timeseriesCollectionName is linked to from these pageType on topics or content nodes
		// contentItems
		// contentHighlightedContent
		var data timeseriesResponse

		if err := json.Unmarshal(fixedJSON, &data); err != nil {
			unmarshalFail(fullURI, err, 11)
		}
		payload, err := json.Marshal(data)
		checkMarshaling(fullURI, err, 11, &payload, &fixedJSON, timeseriesCollectionName)

		title = extractTitle(data.Description.Title)
		description = extractDescription(data.Description.MetaDescription)

		saveContentPageToCollection(timeseriesJsFile, &timeseriesCount, collectionName, bodyTextCopy, shortURI)
		if cfg.FullDepth {
			uriList = getURIListFromTimeseries(fullURI, &data)
		}

	case chartCollectionName:
		// chartCollectionName is linked to from content nodes
		var data chartResponse

		if err := json.Unmarshal(fixedJSON, &data); err != nil {
			unmarshalFail(fullURI, err, 12)
		}
		payload, err := json.Marshal(data)
		checkMarshalingDeepEqual(fullURI, err, 12, &payload, &fixedJSON, chartCollectionName)

		// NOTE: this is different to previous pages ..
		title = extractTitle(data.Title)
		description = extractDescription(data.Subtitle)

		saveContentPageToCollection(chartJsFile, &chartCount, collectionName, bodyTextCopy, shortURI)
		// NOTE: this page has no URI links to add to list

	case productPageCollectionName:
		// NOTE: this is an upper level page being linked back up to
		// ( this should probably not be being grabbed .. we shall see if its of use )
		// productPageCollectionName has been linked to from content nodes
		var data DataResponse

		if err := json.Unmarshal(fixedJSON, &data); err != nil {
			unmarshalFail(fullURI, err, 13)
		}
		payload, err := json.Marshal(data)
		checkMarshaling(fullURI, err, 13, &payload, &fixedJSON, productPageCollectionName)

		title = extractTitle(data.Description.Title)
		description = extractDescription(data.Description.MetaDescription)

		saveContentPageToCollection(productPageJsFile, &productPageCount, collectionName, bodyTextCopy, shortURI)
		// NOTE .. do NOT grab URI's from this as its a top level page from where we initially came.

	case tableCollectionName:
		// tableCollectionName is linked to from content nodes
		var data tableResponse

		if err := json.Unmarshal(fixedJSON, &data); err != nil {
			unmarshalFail(fullURI, err, 14)
		}
		payload, err := json.Marshal(data)
		checkMarshaling(fullURI, err, 14, &payload, &fixedJSON, tableCollectionName)

		// NOTE: this is different to previous pages ..
		title = extractTitle(data.Title)
		description = noDescription

		saveContentPageToCollection(tableJsFile, &tableCount, collectionName, bodyTextCopy, shortURI)
		// NOTE: this page has no URI links to add to list

	case equationCollectionName:
		// equationCollectionName is linked to from content nodes
		var data equationResponse

		if err := json.Unmarshal(fixedJSON, &data); err != nil {
			unmarshalFail(fullURI, err, 15)
		}
		payload, err := json.Marshal(data)
		checkMarshaling(fullURI, err, 15, &payload, &fixedJSON, equationCollectionName)

		// NOTE: this is different to previous pages ..
		title = extractTitle(data.Title)
		description = noDescription

		saveContentPageToCollection(equationJsFile, &equationCount, collectionName, bodyTextCopy, shortURI)
		// NOTE: this page has no URI links to add to list

	case imageCollectionName:
		// imageCollectionName is linked to from content nodes
		var data imageResponse

		if err := json.Unmarshal(fixedJSON, &data); err != nil {
			unmarshalFail(fullURI, err, 16)
		}
		payload, err := json.Marshal(data)
		checkMarshaling(fullURI, err, 16, &payload, &fixedJSON, imageCollectionName)

		// NOTE: this is different to previous pages ..
		title = extractTitle(data.Title)
		description = extractDescription(data.Subtitle)

		saveContentPageToCollection(imageJsFile, &imageCount, collectionName, bodyTextCopy, shortURI)
		// NOTE: this page has no URI links to add to list

	case releaseCollectionName:
		// releaseCollectionName is linked to from content nodes
		var data releaseResponse

		if err := json.Unmarshal(fixedJSON, &data); err != nil {
			unmarshalFail(fullURI, err, 17)
		}
		payload, err := json.Marshal(data)
		checkMarshaling(fullURI, err, 17, &payload, &fixedJSON, releaseCollectionName)

		title = extractTitle(data.Description.Title)
		description = extractDescription(data.Description.MetaDescription)

		saveContentPageToCollection(releaseJsFile, &releaseCount, collectionName, bodyTextCopy, shortURI)
		if cfg.FullDepth {
			uriList = getURIListFromRelease(fullURI, &data)
		}

	case listCollectionName:
		// listCollectionName is linked to from content nodes
		var data listResponse

		if err := json.Unmarshal(fixedJSON, &data); err != nil {
			unmarshalFail(fullURI, err, 18)
		}
		payload, err := json.Marshal(data)
		checkMarshalingDeepEqual(fullURI, err, 18, &payload, &fixedJSON, listCollectionName)

		// NOTE: this is different to previous pages ..
		title = noTitle
		description = noDescription

		saveContentPageToCollection(listJsFile, &listCount, collectionName, bodyTextCopy, shortURI)
		// NOTE: this page has no URI links to add to list

	case staticPageCollectionName:
		// staticPageCollectionName is linked to from content nodes
		var data staticPageResponse

		if err := json.Unmarshal(fixedJSON, &data); err != nil {
			unmarshalFail(fullURI, err, 19)
		}
		payload, err := json.Marshal(data)
		checkMarshaling(fullURI, err, 19, &payload, &fixedJSON, staticPageCollectionName)

		title = extractTitle(data.Description.Title)
		description = extractDescription(data.Description.MetaDescription)

		saveContentPageToCollection(staticPageJsFile, &staticPageCount, collectionName, bodyTextCopy, shortURI)
		if cfg.FullDepth {
			uriList = getURIListFromStaticPage(fullURI, &data)
		}

	case staticAdhocCollectionName:
		// staticAdhocCollectionName is linked to from content nodes
		var data staticAdhocResponse

		if err := json.Unmarshal(fixedJSON, &data); err != nil {
			unmarshalFail(fullURI, err, 20)
		}
		payload, err := json.Marshal(data)
		checkMarshaling(fullURI, err, 20, &payload, &fixedJSON, staticAdhocCollectionName)

		title = extractTitle(data.Description.Title)
		description = extractDescription(data.Description.MetaDescription)

		saveContentPageToCollection(staticAdhocJsFile, &staticAdhocCount, collectionName, bodyTextCopy, shortURI)
		if cfg.FullDepth {
			uriList = getURIListFromStaticAdhoc(fullURI, &data)
		}

	case referenceTablesCollectionName:
		// referenceTablesCollectionName is linked to from content nodes
		var data referenceTablesResponse

		if err := json.Unmarshal(fixedJSON, &data); err != nil {
			unmarshalFail(fullURI, err, 21)
		}
		payload, err := json.Marshal(data)
		checkMarshaling(fullURI, err, 21, &payload, &fixedJSON, referenceTablesCollectionName)

		title = extractTitle(data.Description.Title)
		description = extractDescription(data.Description.MetaDescription)

		saveContentPageToCollection(referenceTablesJsFile, &referenceTablesCount, collectionName, bodyTextCopy, shortURI)
		if cfg.FullDepth {
			uriList = getURIListFromReferenceTables(fullURI, &data)
		}

	case compendiumChapterCollectionName:
		// compendiumChapterCollectionName is linked to from content nodes
		var data compendiumChapterResponse

		if err := json.Unmarshal(fixedJSON, &data); err != nil {
			unmarshalFail(fullURI, err, 22)
		}
		payload, err := json.Marshal(data)
		checkMarshaling(fullURI, err, 22, &payload, &fixedJSON, compendiumChapterCollectionName)

		title = extractTitle(data.Description.Title)
		description = extractDescription(data.Description.MetaDescription)

		saveContentPageToCollection(compendiumChapterJsFile, &compendiumChapterCount, collectionName, bodyTextCopy, shortURI)
		if cfg.FullDepth {
			uriList = getURIListFromCompendiumChapter(fullURI, &data)
		}

	case staticLandingPageCollectionName:
		// staticLandingPageCollectionName is linked to from content nodes
		var data staticLandingPageResponse

		if err := json.Unmarshal(fixedJSON, &data); err != nil {
			unmarshalFail(fullURI, err, 23)
		}
		payload, err := json.Marshal(data)
		checkMarshaling(fullURI, err, 23, &payload, &fixedJSON, staticLandingPageCollectionName)

		title = extractTitle(data.Description.Title)
		description = extractDescription(data.Description.MetaDescription)

		saveContentPageToCollection(staticLandingPageJsFile, &staticLandingPageCount, collectionName, bodyTextCopy, shortURI)
		if cfg.FullDepth {
			uriList = getURIListFromStaticLandingPage(fullURI, &data)
		}

	case staticArticleCollectionName:
		// staticArticleCollectionName is linked to from content nodes
		var data staticArticleResponse

		if err := json.Unmarshal(fixedJSON, &data); err != nil {
			unmarshalFail(fullURI, err, 24)
		}
		payload, err := json.Marshal(data)
		checkMarshaling(fullURI, err, 24, &payload, &fixedJSON, staticArticleCollectionName)

		title = extractTitle(data.Description.Title)
		description = extractDescription(data.Description.MetaDescription)

		saveContentPageToCollection(staticArticleJsFile, &staticArticleCount, collectionName, bodyTextCopy, shortURI)
		if cfg.FullDepth {
			uriList = getURIListFromStaticArticle(fullURI, &data)
		}

	case datasetCollectionName:
		// datasetCollectionName is linked to from content nodes
		var data datasetResponse

		if err := json.Unmarshal(fixedJSON, &data); err != nil {
			unmarshalFail(fullURI, err, 25)
		}
		payload, err := json.Marshal(data)
		checkMarshaling(fullURI, err, 25, &payload, &fixedJSON, datasetCollectionName)

		title = extractTitle(data.Description.Title)
		description = extractDescription(data.Description.MetaDescription)

		saveContentPageToCollection(datasetJsFile, &datasetCount, collectionName, bodyTextCopy, shortURI)
		if cfg.FullDepth {
			uriList = getURIListFromDataset(fullURI, &data)
		}

	case timeseriesDatasetCollectionName:
		// timeseriesDatasetCollectionName is linked to from content nodes
		var data timeseriesDatasetResponse

		if err := json.Unmarshal(fixedJSON, &data); err != nil {
			unmarshalFail(fullURI, err, 26)
		}
		payload, err := json.Marshal(data)
		checkMarshaling(fullURI, err, 26, &payload, &fixedJSON, timeseriesDatasetCollectionName)

		title = extractTitle(data.Description.Title)
		description = extractDescription(data.Description.MetaDescription)

		saveContentPageToCollection(timeseriesDatasetJsFile, &timeseriesDatasetCount, collectionName, bodyTextCopy, shortURI)
		if cfg.FullDepth {
			uriList = getURIListFromTimeseriesDataset(fullURI, &data)
		}

	case taxonomyLandingPageCollectionName:
		// NOTE: this is an upper level page being linked back up to !!!
		// ( this should probably not be being grabbed .. we shall see if its of use )
		// taxonomyLandingPageCollectionName has been linked to from content nodes
		var data DataResponse

		if err := json.Unmarshal(fixedJSON, &data); err != nil {
			unmarshalFail(fullURI, err, 27)
		}
		payload, err := json.Marshal(data)
		checkMarshaling(fullURI, err, 27, &payload, &fixedJSON, taxonomyLandingPageCollectionName)

		title = extractTitle(data.Description.Title)
		description = extractDescription(data.Description.MetaDescription)

		saveContentPageToCollection(taxonomyLandingPageJsFile, &taxonomyLandingPageCount, collectionName, bodyTextCopy, shortURI)
		// NOTE .. do NOT grab URI's from this as its a top level page from where we initially came.

	default:
		fmt.Printf("Unknown page Type ..\n")
		fmt.Printf("shape: %s\n", *shape.Type)
		fmt.Printf("URI: %s\n", fullURI)

		_, err = fmt.Fprintf(bodyTextFile, "Unknown JSON body:\n")
		check(err)
		_, err = bodyTextFile.Write(bodyTextCopy)
		check(err)
		_, err = fmt.Fprintf(bodyTextFile, "\n")
		check(err)

		os.Exit(82)

		// NOTE:
		//
		// home_page : whose uri is "/" .. this would need custom processing to explicitly add sub uri's similar to
		//             what i'm doing in the fabricate function.
		//
		// taxonomy_landing_page : is the the first level down from 'home_page'
		//
		// product_page : is the level down from 'taxonomy_landing_page'
		//
	}

	// good 200 response, save page for later
	appearanceInfo[shortURI]++
	listOfDuplicateInfo = append(listOfDuplicateInfo, duplicateInfo{
		id:        parentTopicNumber,
		pageType:  pType,
		parentURI: parentURI,
		shortURI:  shortURI,
	})
	listOfPageData = append(listOfPageData, pageData{
		id:              parentTopicNumber,
		subSectionIndex: index,
		pageType:        pType,
		uriStatus:       pType,
		shortURI:        shortURI,
		parentURI:       parentURI,
		fixedPayload:    []byte{},
		title:           title,
		description:     description,
	})

	depth++
	if depth > maxDepth {
		maxDepth = depth
	}

	defer func() {
		depth--
	}()

	currentTime := time.Now()
	fmt.Println(currentTime.Format("2006.01.02 15:04:05"))

	fmt.Printf("Index: %d   Depth: %d - %v - %v : %v\n", indexNumber, depth, parentTopicNumber, pType, shortURI)

	if len(uriList) > 0 {
		// iterate over list and call getPageDataRetry
		// doing this should only increase the size of the collections other that topic or content
		// and have no impact on graphs.
		// it will also probably increase the number of duplicates.

		for _, subURI := range uriList {
			// Recurse thru sub URI's, if not already seen (or some other exclusion applies)

			if strings.Contains(subURI, "http://www.ons.gov.uk") {
				fmt.Printf("WARNING: bad link to site using only HTTP and NOT HTTPS: %s\n", subURI)
			}

			// some of the URI links have the 'ons' site in them which we don't want, so remove if present:
			subURI = strings.ReplaceAll(subURI, "https://www.ons.gov.uk", "")
			subURI = strings.ReplaceAll(subURI, "http://www.ons.gov.uk", "")

			if strings.Contains(subURI, "https://") || strings.Contains(subURI, "http://") {
				fmt.Printf("External site: %s\n", subURI)
				continue
			}
			if subURI[0] != '/' {
				fmt.Printf("Adding missing forward slash to: %s\n", subURI)
				// In at least one place on ONS site a URI was missing a forward slash as the first character
				// and that breaks the attempt to open the URI in the code, so we add the missing '/'
				subURI = "/" + subURI
			}
			if contentDuplicateCheck[subURI] > 0 {
				moreThanOnceURI[subURI]++
				//					fmt.Printf("Already processed: %s\n", subURI)
				continue
			}
			if strings.HasPrefix(subURI, "/ons/external-links/") {
				fmt.Printf("A URI to external site: %s\n%s\n", subURI, fullURI)
				continue
			}
			if strings.HasPrefix(subURI, "/ons/rel/") {
				fmt.Printf("A URI to external site: %s\n%s\n", subURI, fullURI)
				continue
			}
			if strings.HasSuffix(subURI, ".doc") {
				fmt.Printf("A URI to .doc file: %s\n%s\n", subURI, fullURI)
				continue
			}
			if strings.HasSuffix(subURI, "/index.html") {
				fmt.Printf("A URI to /index.html: %s\n%s\n", subURI, fullURI)
				continue
			}
			if strings.Contains(subURI, "#") {
				hashURI[subURI]++
				fmt.Printf("A URI with a '#': %s\n%s\n", subURI, fullURI)
				continue
			}
			if strings.Contains(subURI, "?") {
				questionURI[subURI]++
				fmt.Printf("A URI with a '?': %s\n%s\n", subURI, fullURI)
				continue
			}
			parts := strings.Split(subURI, "/")
			last := parts[len(parts)-1]
			var versionFound bool
			if cfg.SkipVersions {
				if len(last) > 1 {
					if last[0] == 'v' && (last[1] >= '0' && last[1] <= '9') {
						// we found what looks like a version number on the end of the URI path
						versionFound = true
					}
				}
			}
			if versionFound {
				skippedVersionURI[subURI]++
				fmt.Printf("Skipping URI with version on end: %s\n", subURI)
				continue
			}
			// do the recursion ..
			getPageDataRetry(0, subURI, 0, pType, shortURI, bodyTextFile, checkFile)
		}
	}

	return response.StatusCode, collectionName
}

var skippedVersionURI = make(map[string]int) // key: shortURI, value: count unique URI's with version number that has been skipped
var hashURI = make(map[string]int)           // key: shortURI, value: count unique URI's with HASH that has been skipped
var questionURI = make(map[string]int)       // key: shortURI, value: count unique URI's with question mark that has been skipped
var moreThanOnceURI = make(map[string]int)   // key: shortURI, value: count URI's seen more than once

func getPageDataRetry(index int, shortURI string, parentTopicNumber int, pType allowedPageType, parentFullURI string, bodyTextFile io.Writer, checkFile io.Writer) (bool, string) {
	var backOff int = 71
	var status int
	var lType string

	for {
		status, lType = getPageData(shortURI, parentTopicNumber, pType, parentFullURI, index)
		if status == 200 {
			return true, lType
		}
		if status == 404 || status == 503 {
			break
		}
		// got error 429 due to making too many requests in a short period of time
		fmt.Printf("backing Off for: %v\n", backOff)
		for delay := 0; delay < backOff; delay++ {
			time.Sleep(1 * time.Second)
			if delay%3 == 0 {
				fmt.Printf("Index: %d   Seconds remaining: %d\n", indexNumber, backOff-delay)
			}
		}
		backOff += 60
		if backOff > 200 {
			// probably a broken URIm but go try without /data on the end ..
			status = 404
			break
		}
	}
	if status == 404 {
		// try reading page without data on the end ..
		noDataURI := "https://www.ons.gov.uk" + shortURI
		// noDataURI := "https://www.production.onsdigital.co.uk" + shortURI

		var response *http.Response
		var err error
		var attempts int

		fmt.Printf("\nGetting /data failed, trying without /data to look for 'redirect'\n")
		for {
			attemptedGetCount++
			if cfg.PlayNice {
				// a little delay to play nice with ONS site and 'hopefully' not have cloudflare 'reset' the connection
				time.Sleep(100 * time.Millisecond)
			}
			response, err = http.Get(noDataURI)
			if err != nil {
				fmt.Println(err)
				fmt.Printf("getPageDataRetry: http.Get(noDataURI) failed\n")
				fmt.Printf("We now fabricate the response code to a 429 to instigate a retry after a delay 3\n")
				doAndShowDelay71()
			} else {
				if response.StatusCode == 429 {
					response.Body.Close()
					doAndShowDelay71()
				} else {
					break
				}
			}
			attempts++
			if attempts >= 3 {
				// Possible problems are:
				// 1. URI on ONS is broke
				// 2. ONS site is down
				// 3. Network connection to ONS is down
				// SO, give up on this URI ..
				fmt.Printf("URI does not exist:  %v\n", shortURI)
				appearanceInfo[shortURI]++
				listOfDuplicateInfo = append(listOfDuplicateInfo, duplicateInfo{
					id:        parentTopicNumber,
					pageType:  pType,
					parentURI: parentFullURI,
					shortURI:  shortURI,
				})
				listOfPageData = append(listOfPageData, pageData{
					id:              parentTopicNumber,
					subSectionIndex: index,
					pageType:        pType,
					uriStatus:       pageBroken,
					shortURI:        shortURI,
					parentURI:       parentFullURI,
					fixedPayload:    []byte{},
				})

				return false, ""
			}
		}

		defer response.Body.Close()

		redirectedURI := response.Request.URL.Path
		if strings.HasPrefix(redirectedURI, "//") {
			fmt.Printf("Double slash    : %s\n", redirectedURI)
			// Remove the first slash as it messes up the usage of the link
			// (the double slash signifies the link uses the same http or https protocol as whatever was being used)
			redirectedURI = redirectedURI[1:]
		}
		fmt.Printf("failed to get URI: %v\n", shortURI)

		if shortURI != redirectedURI {
			fmt.Printf("redirected    URI: %v\n", redirectedURI)
			// we have a re-direction, so lets try that with /data
			backOff = 71
			for {
				status, lType = getPageData(redirectedURI, parentTopicNumber, pType, parentFullURI, index)
				if status == 200 {
					// redirect worked, and page data was saved in the call to getPageData()
					fmt.Printf("redirect worked OK\n")
					return true, lType
				}
				if status == 404 || status == 503 {
					break
				}
				// got error 429 due to making too many requests in a short period of time
				fmt.Printf("backing Off for: %v\n", backOff)
				for delay := 0; delay < backOff; delay++ {
					time.Sleep(1 * time.Second)
					if delay%3 == 0 {
						fmt.Printf("Index: %d   Seconds remaining: %d\n", indexNumber, backOff-delay)
					}
				}
				backOff += 60
				if backOff > 200 {
					// probably a broken URI,so give up
					status = 404
					break
				}
			}
		}
		// NOTE: a true 503 error will already have been recorded in getPageData()
		if status == 404 {
			// if we got a 404 then shortURI or redirectedURI is definitely broke

			fmt.Printf("\nERROR on ONS website /data field: %v\n\n", status)
			if shortURI != redirectedURI {
				fmt.Printf("redirected URI does not exist:  %v\n", redirectedURI)
			}
			fmt.Printf("URI does not exist:  %v\n", shortURI)
			appearanceInfo[shortURI]++
			listOfDuplicateInfo = append(listOfDuplicateInfo, duplicateInfo{
				id:        parentTopicNumber,
				pageType:  pType,
				parentURI: parentFullURI,
				shortURI:  shortURI,
			})
			listOfPageData = append(listOfPageData, pageData{
				id:              parentTopicNumber,
				subSectionIndex: index,
				pageType:        pType,
				uriStatus:       pageBroken,
				shortURI:        shortURI,
				parentURI:       parentFullURI,
				fixedPayload:    []byte{},
			})
		}
	}

	return false, ""
}

/*
root node
node
terminal node
*/

func getTerminationNodeData(data *DataResponse, parentTopicNumber int, parentFullURI string) contentInfo {
	var info contentInfo = contentNone

	if !cfg.ScrapeContent {
		// skip looking for content
		return contentUnknown
	}

	// read any child 'Items                     (Timeseries) links' and save their page /data
	if data.Items != nil {
		if len(*data.Items) > 0 {
			fmt.Printf("Getting: Items\n")
			for index, link := range *data.Items {
				valid, lType := getPageDataRetry(index, *link.URI, parentTopicNumber, pageContentItems, parentFullURI, bodyTextFile, checkFile)
				if valid {
					info = contentExists
					(*data.Items)[index].valid = true
					(*data.Items)[index].linkType = lType
				}
			}
		}
	}

	// read any child 'Datasets                  (Static datasets) links' and save their page /data
	if data.Datasets != nil {
		if len(*data.Datasets) > 0 {
			fmt.Printf("Getting: Datasets\n")
			for index, link := range *data.Datasets {
				valid, lType := getPageDataRetry(index, *link.URI, parentTopicNumber, pageContentDatasets, parentFullURI, bodyTextFile, checkFile)
				if valid {
					info = contentExists
					(*data.Datasets)[index].valid = true
					(*data.Datasets)[index].linkType = lType
				}
			}
		}
	}

	// read any child 'StatsBulletins            (Bulletins) links' and save their page /data
	if data.StatsBulletins != nil {
		if len(*data.StatsBulletins) > 0 {
			fmt.Printf("Getting: StatsBulletins\n")
			for index, link := range *data.StatsBulletins {
				valid, lType := getPageDataRetry(index, *link.URI, parentTopicNumber, pageContentStatsBulletins, parentFullURI, bodyTextFile, checkFile)
				if valid {
					info = contentExists
					(*data.StatsBulletins)[index].valid = true
					(*data.StatsBulletins)[index].linkType = lType
				}
			}
		}
	}

	// read any child 'RelatedArticles           (Articles) links' and save their page /data
	if data.RelatedArticles != nil {
		if len(*data.RelatedArticles) > 0 {
			fmt.Printf("Getting: RelatedArticles\n")
			for index, link := range *data.RelatedArticles {
				valid, lType := getPageDataRetry(index, *link.URI, parentTopicNumber, pageContentRelatedArticles, parentFullURI, bodyTextFile, checkFile)
				if valid {
					info = contentExists
					(*data.RelatedArticles)[index].valid = true
					(*data.RelatedArticles)[index].linkType = lType
				}
			}
		}
	}

	// read any child 'RelatedMethodology        (Methodologies) links' and save their page /data
	if data.RelatedMethodology != nil {
		if len(*data.RelatedMethodology) > 0 {
			fmt.Printf("Getting: RelatedMethodology\n")
			for index, link := range *data.RelatedMethodology {
				valid, lType := getPageDataRetry(index, *link.URI, parentTopicNumber, pageContentRelatedMethodology, parentFullURI, bodyTextFile, checkFile)
				if valid {
					info = contentExists
					(*data.RelatedMethodology)[index].valid = true
					(*data.RelatedMethodology)[index].linkType = lType
				}
			}
		}
	}

	// read any child 'RelatedMethodologyArticle (Methodology_articles) links' and save their page /data
	if data.RelatedMethodologyArticle != nil {
		if len(*data.RelatedMethodologyArticle) > 0 {
			fmt.Printf("Getting: RelatedMethodologyArticle\n")
			for index, link := range *data.RelatedMethodologyArticle {
				valid, lType := getPageDataRetry(index, *link.URI, parentTopicNumber, pageContentRelatedMethodologyArticle, parentFullURI, bodyTextFile, checkFile)
				if valid {
					info = contentExists
					(*data.RelatedMethodologyArticle)[index].valid = true
					(*data.RelatedMethodologyArticle)[index].linkType = lType
				}
			}
		}
	}

	// read any child 'HighlightedContent        (Spotlight) links' and save their page /data
	if data.HighlightedContent != nil {
		if len(*data.HighlightedContent) > 0 {
			fmt.Printf("Getting: HighlightedContent\n")
			for index, link := range *data.HighlightedContent {
				valid, lType := getPageDataRetry(index, *link.URI, parentTopicNumber, pageContentHighlightedContent, parentFullURI, bodyTextFile, checkFile)
				if valid {
					info = contentExists
					(*data.HighlightedContent)[index].valid = true
					(*data.HighlightedContent)[index].linkType = lType
				}
			}
		}
	}

	return info
}

func getNodeData(data *DataResponse, parentTopicNumber int, parentFullURI string, bodyTextFile io.Writer, checkFile io.Writer) (anyValid bool) {
	// read any child 'highlighted links' and save their page /data
	if data.HighlightedLinks != nil {
		if len(*data.HighlightedLinks) > 0 {
			fmt.Printf("Getting: HighlightedLinks\n")
			for index, link := range *data.HighlightedLinks {
				// process tha URI as a 'pageContentHighlightedContent', instead of a 'pageTopicHighlightedLinks'
				// to push it onto a content page
				valid, lType := getPageDataRetry(index, *link.URI, parentTopicNumber, pageContentHighlightedContent, parentFullURI, bodyTextFile, checkFile)
				if valid {
					(*data.HighlightedLinks)[index].valid = true
					(*data.HighlightedLinks)[index].linkType = lType
					anyValid = true
				}
			}
		}
	}

	return anyValid
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func populateTopicAndContentStructs(topics []TopicResponseStore, content []ContentResponse) {

	fmt.Printf("\nindexNumber: %d\n", indexNumber)

	indexNamesLen := len(indexNames)
	fmt.Printf("\nLength of indexNames: %d\n", indexNamesLen)

	var lastPart []string
	var idRef []string
	idRef = append(idRef, "dummy") // offset the subsequent additions so that they can be indexed from 1 with 'id' number
	for i := 1; i <= indexNamesLen; i++ {
		parts := strings.Split(indexNames[i], "/")
		if len(parts) == 0 {
			fmt.Printf("what - zero length")
		}
		last := parts[len(parts)-1]
		lastPart = append(lastPart, last)
		idRef = append(idRef, last)
	}
	fmt.Printf("length of lastPart: %d\n", len(lastPart))
	sort.Strings(lastPart)

	uniqueParts := uniqueListString(lastPart)
	sort.Strings(uniqueParts)
	// uniqueParts[4] = uniqueParts[4] + "_-_test" // put this in to check that test below works

	if len(lastPart) != len(uniqueParts) {
		fmt.Printf("SHOW STOPPER - duplicates exist in the end part of indexNames, inspect following list:\n")
		for _, part := range lastPart {
			fmt.Printf("%s\n", part)
		}
		os.Exit(90)
	}
	// check values are the same
	for i, part := range lastPart {
		if part != uniqueParts[i] {
			fmt.Printf("SHOW STOPPER - non unique end part of Href for use as 'id' in topic and content collections:\n")
			fmt.Printf("%d : %s - %s\n", i, part, uniqueParts[i])
			os.Exit(91)
		}
	}

	for id := 1; id <= indexNumber; id++ {
		var next TopicStore
		// pre-assign next struct to allow for cross assignment of SubtopicsIds
		topics[id].Next = &next
	}

	var topicPageCount int
	var contentPageCount int

	var pageType = make([]allowedPageType, indexNumber+1)

	var maxContentItems int

	for id := 1; id <= indexNumber; id++ {

		var contentState string // this acts as a flag, as well as holding value
		var spotlight []TypeLinkObject = []TypeLinkObject{}
		var articles []TypeLinkObject = []TypeLinkObject{}
		var bulletins []TypeLinkObject = []TypeLinkObject{}
		var methodologies []TypeLinkObject = []TypeLinkObject{}
		var methodologyArticles []TypeLinkObject = []TypeLinkObject{}
		var staticDatasets []TypeLinkObject = []TypeLinkObject{}
		var timeseries []TypeLinkObject = []TypeLinkObject{}

		for _, data := range listOfPageData {
			if data.id == id {
				switch data.pageType {
				case pageBroken:
					// this should not happen .. (it's not been seen on a full site scan)
					fmt.Printf("oops: pageBroken\n")
				case pageTopic:
					if pageType[id] != pageTopicAndContent {
						if pageType[id] == pageContent {
							// add topic to content
							pageType[id] = pageTopicAndContent
						} else {
							pageType[id] = pageTopic
						}
					}
					topicPageCount++
					parentID := data.subSectionIndex
					// idAndName := strconv.Itoa(data.id) + " - " + data.shortURI
					idAndName := idRef[id]

					if data.uriStatus == pageTopicBroken {
						pageType[id] = pageTopicBroken
						// The URI that was trying to be viewed does not exist.
						// Therefore we don't have any info on it to put into the topic database.
						// All we can do to indicate a broken link is assign a sub-topics id that has nothing in it in
						// the mongo database (a blank place holder).
						// NOTE: as of 30th Dec' 2020 the dp-topic-api subtopics endpoint can't indicate to the caller
						//       if a subtopic link is broken if it is in a list of 2 or more subtopics and at least one
						//       of them is OK .. check that this is still the case !!!
						// !!! it might be that we don't want to carry broken links forward and code in this if block
						// needs to be removed ..
						if topics[parentID].Next.SubtopicIds == nil {
							topics[parentID].Next.SubtopicIds = &[]string{idAndName}
						} else if !stringInSlice(idAndName, *topics[parentID].Next.SubtopicIds) {
							*topics[parentID].Next.SubtopicIds = append(*topics[parentID].Next.SubtopicIds, idAndName)
						}
					} else {
						topics[id].ID = idAndName
						topics[id].Next.ID = idAndName

						var savedPageData DataResponse

						// Unmarshal body bytes to model
						if err := json.Unmarshal(data.fixedPayload, &savedPageData); err != nil {
							fmt.Println(err)
							fmt.Printf("populateTopicAndContentStructs: json.Unmarshal for pageTopic failed\n")
							os.Exit(10)
						}

						// get & save:
						if savedPageData.Description.MetaDescription != nil {
							topics[id].Next.Description = *savedPageData.Description.MetaDescription
						} else {
							topics[id].Next.Description = noDescription
						}
						topics[id].Next.Title = *savedPageData.Description.Title

						if savedPageData.Description.Keywords != nil {
							for _, k := range *savedPageData.Description.Keywords {
								if topics[id].Next.Keywords == nil {
									topics[id].Next.Keywords = &[]string{k}
								} else {
									*topics[id].Next.Keywords = append(*topics[id].Next.Keywords, k)
								}
							}
						}

						topics[id].Next.State = "published"
					}

				case pageTopicBroken:
					// this should not happen .. (it's not been seen on a full site scan)
					fmt.Printf("oops: pageTopicBroken\n")

				case pageTopicSubtopicID:
					// Add topic node id into parent SubtopicsIds list
					parentID := data.subSectionIndex
					// idAndName := strconv.Itoa(data.id) + " - " + data.shortURI
					idAndName := idRef[id]
					if topics[parentID].Next.SubtopicIds == nil {
						topics[parentID].Next.SubtopicIds = &[]string{idAndName}
					} else if !stringInSlice(idAndName, *topics[parentID].Next.SubtopicIds) {
						*topics[parentID].Next.SubtopicIds = append(*topics[parentID].Next.SubtopicIds, idAndName)
					}

				case pageContent:
					contentState = "published"

					// Add topic termination node id into parent SubtopicsIds list
					parentID := data.subSectionIndex
					// idAndName := strconv.Itoa(data.id) + " - " + data.shortURI
					idAndName := idRef[id]
					if topics[parentID].Next.SubtopicIds == nil {
						topics[parentID].Next.SubtopicIds = &[]string{idAndName}
					} else if !stringInSlice(idAndName, *topics[parentID].Next.SubtopicIds) {
						*topics[parentID].Next.SubtopicIds = append(*topics[parentID].Next.SubtopicIds, idAndName)
					}

					if pageType[id] != pageTopicAndContent {
						if pageType[id] == pageTopic {
							// add content to topic
							pageType[id] = pageTopicAndContent
						} else {
							pageType[id] = pageContent
						}
					}

					contentPageCount++
					topics[id].ID = idAndName
					topics[id].Next.ID = idAndName

					var savedPageData DataResponse

					// Unmarshal body bytes to model
					if err := json.Unmarshal(data.fixedPayload, &savedPageData); err != nil {
						fmt.Println(err)
						fmt.Printf("populateTopicAndContentStructs: json.Unmarshal for pageContent failed\n")
						os.Exit(11)
					}

					// get & save:
					if savedPageData.Description.MetaDescription != nil {
						topics[id].Next.Description = *savedPageData.Description.MetaDescription
					} else {
						topics[id].Next.Description = noDescription
					}
					topics[id].Next.Title = *savedPageData.Description.Title

					if savedPageData.Description.Keywords != nil {
						for _, k := range *savedPageData.Description.Keywords {
							if topics[id].Next.Keywords == nil {
								topics[id].Next.Keywords = &[]string{k}
							} else {
								*topics[id].Next.Keywords = append(*topics[id].Next.Keywords, k)
							}
						}
					}

					topics[id].Next.State = "published"

				case pageContentItems, // Content Timeseries
					pageContentDatasets,                  // Content Static datasets
					pageContentStatsBulletins,            // Content Bulletins
					pageContentRelatedArticles,           // Content Articles
					pageContentRelatedMethodology,        // Content Methodologies
					pageContentRelatedMethodologyArticle, // Content Methodology_articles
					pageContentHighlightedContent:        // Content Spotlight
					var item TypeLinkObject
					if data.uriStatus == pageBroken {
						item.Title = "** broken **" + data.title
					} else {
						item.Title = data.title
					}
					contentState = "published"
					item.HRef = data.shortURI

					switch data.pageType {
					case pageContentItems: // Content Timeseries
						timeseries = append(timeseries, item)
					case pageContentDatasets: // Content Static datasets
						staticDatasets = append(staticDatasets, item)
					case pageContentStatsBulletins: // Content Bulletins
						bulletins = append(bulletins, item)
					case pageContentRelatedArticles: // Content Articles
						articles = append(articles, item)
					case pageContentRelatedMethodology: // Content Methodologies
						methodologies = append(methodologies, item)
					case pageContentRelatedMethodologyArticle: // Content Methodology_articles
						methodologyArticles = append(methodologyArticles, item)
					case pageContentHighlightedContent: // Content Spotlight
						spotlight = append(spotlight, item)
					}
				}
			}
		}
		if contentState != "" {
			// content[id].ID = strconv.Itoa(id) + " - " + indexNames[id]
			content[id].ID = idRef[id]

			var totalItems int
			var next Content
			// pre-assign next struct to allow for cross assignment of SubtopicsIds
			content[id].Next = &next

			content[id].Next.State = contentState

			if len(spotlight) > 0 {
				totalItems += len(spotlight)
				content[id].Next.Spotlight = &spotlight
			}
			if len(articles) > 0 {
				totalItems += len(articles)
				content[id].Next.Articles = &articles
			}
			if len(bulletins) > 0 {
				totalItems += len(bulletins)
				content[id].Next.Bulletins = &bulletins
			}
			if len(methodologies) > 0 {
				totalItems += len(methodologies)
				content[id].Next.Methodologies = &methodologies
			}
			if len(methodologyArticles) > 0 {
				totalItems += len(methodologyArticles)
				content[id].Next.MethodologyArticles = &methodologyArticles
			}
			if len(staticDatasets) > 0 {
				totalItems += len(staticDatasets)
				content[id].Next.StaticDatasets = &staticDatasets
			}
			if len(timeseries) > 0 {
				totalItems += len(timeseries)
				content[id].Next.Timeseries = &timeseries
			}

			if totalItems > maxContentItems {
				maxContentItems = totalItems
			}

			content[id].Current = content[id].Next
		}
	}

	fmt.Printf("\nmaxContentItems: %d\n", maxContentItems)

	var baseURI string = "http://localhost:25300/topics/"

	// assign topic Links ..
	for id := 1; id <= indexNumber; id++ {
		var self LinkObject
		var subtopics LinkObject
		var content LinkObject
		var topicLinks TopicLinks

		idAndName := idRef[id]

		self.HRef = baseURI + idAndName
		self.ID = idAndName

		topicLinks.Self = &self

		if pageType[id] == pageTopic || pageType[id] == pageTopicAndContent {
			if topics[id].Next.SubtopicIds != nil {
				if len(*topics[id].Next.SubtopicIds) > 0 {
					subtopics.HRef = baseURI + idAndName + "/subtopics"
					topicLinks.Subtopics = &subtopics
				}
			}
		}

		if pageType[id] == pageContent || pageType[id] == pageTopicAndContent {
			content.HRef = baseURI + idAndName + "/content"
			topicLinks.Content = &content
		}

		if pageType[id] != pageTopicBroken {
			topics[id].Next.Links = &topicLinks
			topics[id].Current = topics[id].Next
		} else {
			topics[id].Next = nil
		}
	}

	fmt.Printf("\ntopicPageCount: %d\n", topicPageCount)
	fmt.Printf("\ncontentPageCount: %d\n", contentPageCount)
}

func uniqueListString(stringSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

var topicsDbName = "topics"
var topicsDbCollection = "topics"
var initDir = "mongo-init-scripts"

func createTopicJsScript(topics []TopicResponseStore) {
	// do the topic database ..
	topicsJsFile, err := os.Create(initDir + "/" + topicsDbName + "-init.js")
	check(err)
	defer topicsJsFile.Close()

	line1 := "db = db.getSiblingDB('" + topicsDbName + "')\n"
	line2 := "db." + topicsDbCollection + ".remove({})\n"

	_, err = fmt.Fprint(topicsJsFile, line1)
	check(err)
	_, err = fmt.Fprint(topicsJsFile, line2)
	check(err)

	// write each document (topic node)
	for id := 1; id <= indexNumber; id++ {
		if topics[id].Next != nil {
			//			payload, err := json.Marshal(topics[id])
			payload, err := json.MarshalIndent(topics[id], "", "    ")
			if err != nil {
				fmt.Println(err)
				fmt.Printf("createTopicJsScript: json.Marshal failed\n")
				os.Exit(14)
			}

			_, err = fmt.Fprintf(topicsJsFile, "db."+topicsDbCollection+".insertOne(")
			check(err)
			// Must 'write' the 'payload', as converting it to a string causes percentage characters in the
			// 'payload' to be missinterpreted.
			_, err = topicsJsFile.Write(payload)
			check(err)
			_, err = fmt.Fprintf(topicsJsFile, ")\n")
			check(err)

		}
	}

	// Add code to read back each document written (for visual inspection)
	_, err = fmt.Fprintf(topicsJsFile, "db."+topicsDbCollection+".find().forEach(function(doc) {\n")
	check(err)
	_, err = fmt.Fprintf(topicsJsFile, "    printjson(doc);\n")
	check(err)
	_, err = fmt.Fprintf(topicsJsFile, "})\n")
	check(err)
}

var contentDbCollection = "content"

func createContentJsScript(content []ContentResponse) {
	// do the content database ..
	contentJsFile, err := os.Create(initDir + "/" + contentDbCollection + "-init.js")
	check(err)
	defer contentJsFile.Close()

	line1 := "db = db.getSiblingDB('" + topicsDbName + "')\n"
	line2 := "db." + contentDbCollection + ".remove({})\n"

	_, err = fmt.Fprint(contentJsFile, line1)
	check(err)
	_, err = fmt.Fprint(contentJsFile, line2)
	check(err)

	// write each document (topic node)
	for id := 1; id <= indexNumber; id++ {
		if content[id].Next != nil {
			//			payload, err := json.Marshal(content[id])
			payload, err := json.MarshalIndent(content[id], "", "    ")
			if err != nil {
				fmt.Println(err)
				fmt.Printf("createContentJsScript: json.Marshal failed\n")
				os.Exit(15)
			}

			_, err = fmt.Fprintf(contentJsFile, "db."+contentDbCollection+".insertOne(")
			check(err)
			// Must 'write' the 'payload', as converting it to a string causes percentage characters in the
			// 'payload' to be missinterpreted.
			_, err = contentJsFile.Write(payload)
			check(err)
			_, err = fmt.Fprintf(contentJsFile, ")\n")
			check(err)
		}
	}

	// Add code to read back each document written (for visual inspection)
	_, err = fmt.Fprintf(contentJsFile, "db."+contentDbCollection+".find().forEach(function(doc) {\n")
	check(err)
	_, err = fmt.Fprintf(contentJsFile, "    printjson(doc);\n")
	check(err)
	_, err = fmt.Fprintf(contentJsFile, "})\n")
	check(err)
}

func createBrokenLinkFile() {
	if listOfPageData != nil {
		fmt.Printf("\nNof listOfPageData: %v\n", len(listOfPageData))
		if len(listOfPageData) > 0 {
			brokenLinkTextFile, err := os.Create(observationsDir + "/broken_links.txt")
			check(err)
			defer brokenLinkTextFile.Close()

			var errorCount int
			fmt.Printf("Showing: listOfPageData\n")
			for _, pagesData := range listOfPageData {
				if pagesData.pageType != pagesData.uriStatus {
					errorCount++
					parentFullURI := pagesData.parentURI
					if parentFullURI[0] == '/' {
						parentFullURI = "https://www.ons.gov.uk" + parentFullURI
					}
					fmt.Printf("Error on page: %v\n    Broken link: ", parentFullURI)

					// save to file
					_, err = fmt.Fprintf(brokenLinkTextFile, "===================\n")
					check(err)
					_, err = fmt.Fprintf(brokenLinkTextFile, "%v - Error on page: %v\n\n", errorCount, parentFullURI)
					check(err)
					_, err = fmt.Fprintf(brokenLinkTextFile, "%v:\n", pageTypeString(pagesData.pageType))
					check(err)
					_, err = fmt.Fprintf(brokenLinkTextFile, "  %v:\n", pagesData.subSectionIndex)
					check(err)
					_, err = fmt.Fprintf(brokenLinkTextFile, "    Broken link: uri: %v\n\n", pagesData.shortURI)
					check(err)
					_, err = fmt.Fprintf(brokenLinkTextFile, "    Broken link: https: %v\n\n", "https://www.ons.gov.uk"+pagesData.shortURI)
					check(err)
				}
			}
		}
	}
}

// iterate through appearanceInfo and where val > 1, use the key to
// search for key in listOfDuplicateInfo and print out the parentURI
func createDuplicatesFile() {
	var nofDuplicates int

	defer func() {
		fmt.Printf("\nNof duplicates: %v\n", nofDuplicates)
	}()

	if appearanceInfo == nil {
		return
	}
	if listOfDuplicateInfo == nil {
		return
	}
	if len(appearanceInfo) == 0 {
		return
	}

	for shortURI := range appearanceInfo {
		if appearanceInfo[shortURI] > 1 {
			nofDuplicates++
		}
	}

	if nofDuplicates == 0 {
		return
	}

	duplicatesLinkTextFile, err := os.Create(observationsDir + "/duplicates_links.txt")
	check(err)
	defer duplicatesLinkTextFile.Close()

	_, err = fmt.Fprintf(duplicatesLinkTextFile, "Nof duplicates: %v\n", nofDuplicates)
	check(err)

	// Create a list of sorted URIs
	URIs := make([]string, 0, len(appearanceInfo))
	for k := range appearanceInfo {
		URIs = append(URIs, k)
	}
	sort.Strings(URIs)

	// Use sorted list of URIs to iterate through 'appearanceInfo' in order so that the output
	// file can be comparred to a previous output file and the differences visually align
	// when opening both files in an application like 'meld'.
	for _, shortURI := range URIs {
		if appearanceInfo[shortURI] > 1 {
			_, err = fmt.Fprintf(duplicatesLinkTextFile, "\nPage: %v\nIs linked to from:\n", shortURI)
			check(err)

			for _, duplicate := range listOfDuplicateInfo {
				if shortURI == duplicate.shortURI {
					_, err = fmt.Fprintf(duplicatesLinkTextFile, "    %v\n", duplicate.parentURI)
					check(err)
				}
			}
		}
	}
}

// create file that contains list of URI's saved when doing deeper scan together with
// the name of the collection that the URI info is storred in - that is the 'type'
// of the page and thus one knows the struct to use to read the URI
func createURICollectionNamesFile() {
	var nofURIs int

	defer func() {
		fmt.Printf("\nNof URI's / keys storred: %v\n", nofURIs)
	}()

	if uriCollectionName == nil {
		return
	}
	if len(uriCollectionName) == 0 {
		return
	}

	nofURIs = len(uriCollectionName)

	if nofURIs == 0 {
		return
	}

	namesTextFile, err := os.Create("mongo-init-scripts/uri_collection_names.txt")
	check(err)
	defer namesTextFile.Close()

	// Create a list of sorted URIs
	URIs := make([]string, 0, nofURIs)
	for k := range uriCollectionName {
		URIs = append(URIs, k)
	}
	sort.Strings(URIs)

	// Use sorted list of URIs to iterate through 'uriCollectionName' in order to
	// save the URI's and their collection name in order of URI's
	for _, shortURI := range URIs {
		_, err = fmt.Fprintf(namesTextFile, "%s,%s\n", shortURI, uriCollectionName[shortURI])
		check(err)
	}
}

func createContentCountsFile() {
	countsTextFile, err := os.Create("mongo-init-scripts/collection_lengths.txt")
	check(err)
	defer countsTextFile.Close()

	_, err = fmt.Fprintf(countsTextFile, "article_download collection quantity: %d\n", articleDownloadCount)
	check(err)
	_, err = fmt.Fprintf(countsTextFile, "article collection quantity: %d\n", articleCount)
	check(err)
	_, err = fmt.Fprintf(countsTextFile, "bulletin collection quantity: %d\n", bulletinCount)
	check(err)
	_, err = fmt.Fprintf(countsTextFile, "chart collection quantity: %d\n", chartCount)
	check(err)
	_, err = fmt.Fprintf(countsTextFile, "compendium_chapter collection quantity: %d\n", compendiumChapterCount)
	check(err)
	_, err = fmt.Fprintf(countsTextFile, "compendium_data collection quantity: %d\n", compendiumDataCount)
	check(err)
	_, err = fmt.Fprintf(countsTextFile, "compendium_landing_page collection quantity: %d\n", compendiumLandingPageCount)
	check(err)
	_, err = fmt.Fprintf(countsTextFile, "dataset_landing_page collection quantity: %d\n", datasetLandingPageCount)
	check(err)
	_, err = fmt.Fprintf(countsTextFile, "dataset collection quantity: %d\n", datasetCount)
	check(err)
	_, err = fmt.Fprintf(countsTextFile, "equation collection quantity: %d\n", equationCount)
	check(err)
	_, err = fmt.Fprintf(countsTextFile, "image collection quantity: %d\n", imageCount)
	check(err)
	_, err = fmt.Fprintf(countsTextFile, "list collection quantity: %d\n", listCount)
	check(err)
	_, err = fmt.Fprintf(countsTextFile, "product_page collection quantity: %d\n", productPageCount)
	check(err)
	_, err = fmt.Fprintf(countsTextFile, "reference_tables collection quantity: %d\n", referenceTablesCount)
	check(err)
	_, err = fmt.Fprintf(countsTextFile, "release collection quantity: %d\n", releaseCount)
	check(err)
	_, err = fmt.Fprintf(countsTextFile, "static_adhoc collection quantity: %d\n", staticAdhocCount)
	check(err)
	_, err = fmt.Fprintf(countsTextFile, "static_article collection quantity: %d\n", staticArticleCount)
	check(err)
	_, err = fmt.Fprintf(countsTextFile, "static_landing_page collection quantity: %d\n", staticLandingPageCount)
	check(err)
	_, err = fmt.Fprintf(countsTextFile, "static_methodology_download collection quantity: %d\n", staticMethodologyDownloadCount)
	check(err)
	_, err = fmt.Fprintf(countsTextFile, "static_methodology collection quantity: %d\n", staticMethodologyCount)
	check(err)
	_, err = fmt.Fprintf(countsTextFile, "static_page collection quantity: %d\n", staticPageCount)
	check(err)
	_, err = fmt.Fprintf(countsTextFile, "static_qmi collection quantity: %d\n", staticQmiCount)
	check(err)
	_, err = fmt.Fprintf(countsTextFile, "table collection quantity: %d\n", tableCount)
	check(err)
	_, err = fmt.Fprintf(countsTextFile, "taxonomy_panding_page collection quantity: %d\n", taxonomyLandingPageCount)
	check(err)
	_, err = fmt.Fprintf(countsTextFile, "timeseries_dataset collection quantity: %d\n", timeseriesDatasetCount)
	check(err)
	_, err = fmt.Fprintf(countsTextFile, "timeseries collection quantity: %d\n", timeseriesCount)
	check(err)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

var rootPage DataResponse

var rootPort string = "64000"
var rootPath string = "/topic_root"
var rootURI string = "http://localhost:" + rootPort // local handler that serves up 'rootPage'

func fabricateRootPage() {
	l1 := "/businessindustryandtrade"
	n1 := 0
	l2 := "/economy"
	n2 := 1
	l3 := "/employmentandlabourmarket"
	n3 := 2
	l4 := "/peoplepopulationandcommunity"
	n4 := 3
	var link []SubLink = []SubLink{
		{URI: &l1,
			Index: &n1},
		{URI: &l2,
			Index: &n2},
		{URI: &l3,
			Index: &n3},
		{URI: &l4,
			Index: &n4},
	}

	rootPage.Sections = &link

	var des Descript
	r := "root page"
	des.MetaDescription = &r
	rp := "The root page"
	des.Title = &rp

	var keys []string = []string{"root 1", "root 2"}
	des.Keywords = &keys

	rootPage.Description = &des
}

func rootPageHandler(w http.ResponseWriter, r *http.Request) {
	// Set headers
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// Marshal provided model
	payload, err := json.Marshal(rootPage)
	check(err)

	// Write payload to body
	_, err = w.Write(payload)
	check(err)
}

var articleJsFile *os.File

const articleCollectionName string = "article"

var articleDownloadJsFile *os.File

const articleDownloadCollectionName string = "article_download"

var bulletinJsFile *os.File

const bulletinCollectionName string = "bulletin"

var compendiumDataJsFile *os.File

const compendiumDataCollectionName string = "compendium_data"

var compendiumLandingPageJsFile *os.File

const compendiumLandingPageCollectionName string = "compendium_landing_page"

var datasetLandingPageJsFile *os.File

const datasetLandingPageCollectionName string = "dataset_landing_page"

var staticMethodologyJsFile *os.File

const staticMethodologyCollectionName string = "static_methodology"

var staticMethodologyDownloadJsFile *os.File

const staticMethodologyDownloadCollectionName string = "static_methodology_download"

var staticQmiJsFile *os.File

const staticQmiCollectionName string = "static_qmi"

var timeseriesJsFile *os.File

const timeseriesCollectionName string = "timeseries"

var chartJsFile *os.File

const chartCollectionName string = "chart"

var productPageJsFile *os.File

const productPageCollectionName string = "product_page"

var tableJsFile *os.File

const tableCollectionName string = "table"

var equationJsFile *os.File

const equationCollectionName string = "equation"

var imageJsFile *os.File

const imageCollectionName string = "image"

var releaseJsFile *os.File

const releaseCollectionName string = "release"

var listJsFile *os.File

const listCollectionName string = "list"

var staticPageJsFile *os.File

const staticPageCollectionName string = "static_page"

var staticAdhocJsFile *os.File

const staticAdhocCollectionName string = "static_adhoc"

var referenceTablesJsFile *os.File

const referenceTablesCollectionName string = "reference_tables"

var compendiumChapterJsFile *os.File

const compendiumChapterCollectionName string = "compendium_chapter"

var staticLandingPageJsFile *os.File

const staticLandingPageCollectionName string = "static_landing_page"

var staticArticleJsFile *os.File

const staticArticleCollectionName string = "static_article"

var datasetJsFile *os.File

const datasetCollectionName string = "dataset"

var timeseriesDatasetJsFile *os.File

const timeseriesDatasetCollectionName string = "timeseries_dataset"

var taxonomyLandingPageJsFile *os.File

const taxonomyLandingPageCollectionName string = "taxonomy_landing_page"

var bodyTextFile *os.File
var checkFile *os.File

func initialiseCollectionDatabase(collectionName string, collectionFile *os.File) {
	line1 := "db = db.getSiblingDB('" + topicsDbName + "')\n"
	line2 := "db." + collectionName + ".remove({})\n"

	_, err := fmt.Fprint(collectionFile, line1)
	check(err)
	_, err = fmt.Fprint(collectionFile, line2)
	check(err)
}

func finaliseCollectionDatabase(collectionName string, collectionFile *os.File) {
	// Add code to read back each document written (for visual inspection)
	// NOTE: these lines in script are commented out to speed the process up for long scripts
	//       they are placed in init script should they need to be uncomented ..
	_, err := fmt.Fprint(collectionFile, "//db."+collectionName+".find().forEach(function(doc) {\n")
	check(err)
	_, err = fmt.Fprint(collectionFile, "//    printjson(doc);\n")
	check(err)
	_, err = fmt.Fprint(collectionFile, "//})\n")
	check(err)
}

var graphDir = "graphviz-files"
var tempDir = "temp"
var observationsDir = "observations"

func ensureDirectoryExists(dirName string) {
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		check(os.Mkdir(dirName, 0700))
	}
}

func main() {
	// Read config
	err := InitConfig()
	if err != nil {
		fmt.Printf("error initialising config\n")
		return
	}

	ensureDirectoryExists(initDir)

	// Create files 'article' content creation file
	articleJsFile, err = os.Create(initDir + "/" + articleCollectionName + "-init.js")
	check(err)
	defer articleJsFile.Close()
	initialiseCollectionDatabase(articleCollectionName, articleJsFile)

	// Create files 'article_download' content creation file
	articleDownloadJsFile, err = os.Create(initDir + "/" + articleDownloadCollectionName + "-init.js")
	check(err)
	defer articleDownloadJsFile.Close()
	initialiseCollectionDatabase(articleDownloadCollectionName, articleDownloadJsFile)

	// Create files 'bulletin' content creation file
	bulletinJsFile, err = os.Create(initDir + "/" + bulletinCollectionName + "-init.js")
	check(err)
	defer bulletinJsFile.Close()
	initialiseCollectionDatabase(bulletinCollectionName, bulletinJsFile)

	// Create files 'compendium_data' content creation file
	compendiumDataJsFile, err = os.Create(initDir + "/" + compendiumDataCollectionName + "-init.js")
	check(err)
	defer compendiumDataJsFile.Close()
	initialiseCollectionDatabase(compendiumDataCollectionName, compendiumDataJsFile)

	// Create files 'compendium_landing_page' content creation file
	compendiumLandingPageJsFile, err = os.Create(initDir + "/" + compendiumLandingPageCollectionName + "-init.js")
	check(err)
	defer compendiumLandingPageJsFile.Close()
	initialiseCollectionDatabase(compendiumLandingPageCollectionName, compendiumLandingPageJsFile)

	// Create files 'dataset_landing_page' content creation file
	datasetLandingPageJsFile, err = os.Create(initDir + "/" + datasetLandingPageCollectionName + "-init.js")
	check(err)
	defer datasetLandingPageJsFile.Close()
	initialiseCollectionDatabase(datasetLandingPageCollectionName, datasetLandingPageJsFile)

	// Create files 'static_methodology' content creation file
	staticMethodologyJsFile, err = os.Create(initDir + "/" + staticMethodologyCollectionName + "-init.js")
	check(err)
	defer staticMethodologyJsFile.Close()
	initialiseCollectionDatabase(staticMethodologyCollectionName, staticMethodologyJsFile)

	// Create files 'static_methodology_download' content creation file
	staticMethodologyDownloadJsFile, err = os.Create(initDir + "/" + staticMethodologyDownloadCollectionName + "-init.js")
	check(err)
	defer staticMethodologyDownloadJsFile.Close()
	initialiseCollectionDatabase(staticMethodologyDownloadCollectionName, staticMethodologyDownloadJsFile)

	// Create files 'static_qmi' content creation file
	staticQmiJsFile, err = os.Create(initDir + "/" + staticQmiCollectionName + "-init.js")
	check(err)
	defer staticQmiJsFile.Close()
	initialiseCollectionDatabase(staticQmiCollectionName, staticQmiJsFile)

	// Create files 'timeseries' content creation file
	timeseriesJsFile, err = os.Create(initDir + "/" + timeseriesCollectionName + "-init.js")
	check(err)
	defer timeseriesJsFile.Close()
	initialiseCollectionDatabase(timeseriesCollectionName, timeseriesJsFile)

	// Create files 'chart' content creation file
	chartJsFile, err = os.Create(initDir + "/" + chartCollectionName + "-init.js")
	check(err)
	defer chartJsFile.Close()
	initialiseCollectionDatabase(chartCollectionName, chartJsFile)

	// Create files 'product_page' content creation file
	productPageJsFile, err = os.Create(initDir + "/" + productPageCollectionName + "-init.js")
	check(err)
	defer productPageJsFile.Close()
	initialiseCollectionDatabase(productPageCollectionName, productPageJsFile)

	// Create files 'table' content creation file
	tableJsFile, err = os.Create(initDir + "/" + tableCollectionName + "-init.js")
	check(err)
	defer tableJsFile.Close()
	initialiseCollectionDatabase(tableCollectionName, tableJsFile)

	// Create files 'equation' content creation file
	equationJsFile, err = os.Create(initDir + "/" + equationCollectionName + "-init.js")
	check(err)
	defer equationJsFile.Close()
	initialiseCollectionDatabase(equationCollectionName, equationJsFile)

	// Create files 'image' content creation file
	imageJsFile, err = os.Create(initDir + "/" + imageCollectionName + "-init.js")
	check(err)
	defer imageJsFile.Close()
	initialiseCollectionDatabase(imageCollectionName, imageJsFile)

	// Create files 'release' content creation file
	releaseJsFile, err = os.Create(initDir + "/" + releaseCollectionName + "-init.js")
	check(err)
	defer releaseJsFile.Close()
	initialiseCollectionDatabase(releaseCollectionName, releaseJsFile)

	// Create files 'list' content creation file
	listJsFile, err = os.Create(initDir + "/" + listCollectionName + "-init.js")
	check(err)
	defer listJsFile.Close()
	initialiseCollectionDatabase(listCollectionName, listJsFile)

	// Create files 'static_page' content creation file
	staticPageJsFile, err = os.Create(initDir + "/" + staticPageCollectionName + "-init.js")
	check(err)
	defer staticPageJsFile.Close()
	initialiseCollectionDatabase(staticPageCollectionName, staticPageJsFile)

	// Create files 'static_adhoc' content creation file
	staticAdhocJsFile, err = os.Create(initDir + "/" + staticAdhocCollectionName + "-init.js")
	check(err)
	defer staticAdhocJsFile.Close()
	initialiseCollectionDatabase(staticAdhocCollectionName, staticAdhocJsFile)

	// Create files 'reference_tables' content creation file
	referenceTablesJsFile, err = os.Create(initDir + "/" + referenceTablesCollectionName + "-init.js")
	check(err)
	defer referenceTablesJsFile.Close()
	initialiseCollectionDatabase(referenceTablesCollectionName, referenceTablesJsFile)

	// Create files 'compendium_chapter' content creation file
	compendiumChapterJsFile, err = os.Create(initDir + "/" + compendiumChapterCollectionName + "-init.js")
	check(err)
	defer compendiumChapterJsFile.Close()
	initialiseCollectionDatabase(compendiumChapterCollectionName, compendiumChapterJsFile)

	// Create files 'static_landing_page' content creation file
	staticLandingPageJsFile, err = os.Create(initDir + "/" + staticLandingPageCollectionName + "-init.js")
	check(err)
	defer staticLandingPageJsFile.Close()
	initialiseCollectionDatabase(staticLandingPageCollectionName, staticLandingPageJsFile)

	// Create files 'static_article' content creation file
	staticArticleJsFile, err = os.Create(initDir + "/" + staticArticleCollectionName + "-init.js")
	check(err)
	defer staticArticleJsFile.Close()
	initialiseCollectionDatabase(staticArticleCollectionName, staticArticleJsFile)

	// Create files 'dataset' content creation file
	datasetJsFile, err = os.Create(initDir + "/" + datasetCollectionName + "-init.js")
	check(err)
	defer datasetJsFile.Close()
	initialiseCollectionDatabase(datasetCollectionName, datasetJsFile)

	// Create files 'timeseries_dataset' content creation file
	timeseriesDatasetJsFile, err = os.Create(initDir + "/" + timeseriesDatasetCollectionName + "-init.js")
	check(err)
	defer timeseriesDatasetJsFile.Close()
	initialiseCollectionDatabase(timeseriesDatasetCollectionName, timeseriesDatasetJsFile)

	// Create files 'taxonomy_landing_page' content creation file
	taxonomyLandingPageJsFile, err = os.Create(initDir + "/" + taxonomyLandingPageCollectionName + "-init.js")
	check(err)
	defer taxonomyLandingPageJsFile.Close()
	initialiseCollectionDatabase(taxonomyLandingPageCollectionName, taxonomyLandingPageJsFile)

	ensureDirectoryExists(graphDir)

	gvFile := graphDir + "/t.gv"
	if cfg.GraphAllContent {
		if cfg.ColourContent {
			gvFile = graphDir + "/t-big-colour.gv"
		} else {
			gvFile = graphDir + "/t-big.gv"
		}
	}
	graphVizFile, err := os.Create(gvFile)
	check(err)
	defer graphVizFile.Close()

	ensureDirectoryExists(tempDir)

	// These files are saved for visual comparison when a structure decode and encode differ.
	// Open both files in vscode, right click in them and select 'Format Document' to expand the json,
	// save each expanded .json file and then do a visual diff between them with meld.
	// It is recommended that you use meld because some timeseries can be over 12,000 lines long.
	bodyTextFile, err = os.Create(tempDir + "/bodyText_all.json")
	check(err)
	defer bodyTextFile.Close()

	checkFile, err = os.Create(tempDir + "/bodyText_all_processed.json")
	check(err)
	defer checkFile.Close()

	ensureDirectoryExists(observationsDir)

	// Open the graph
	_, err = fmt.Fprintf(graphVizFile, "digraph G {\n\n    rankdir=LR\n    ranksep=2.7\n\n")
	check(err)

	fabricateRootPage()

	http.HandleFunc("/", rootPageHandler)
	go func() {
		err = http.ListenAndServe(":"+rootPort, nil)
		check(err)
	}()

	// give server a little time to start before accessing it ..
	time.Sleep(1 * time.Second)

	// iterate and recurse through ONS site starting at specified: rootPath ..
	getPage(1, graphVizFile, "", rootPath)

	// Close the whole graph:
	_, err = fmt.Fprintf(graphVizFile, "}\n")
	check(err)

	fmt.Printf("\nTotal GOOD pages: %d\n", pageCount)

	fmt.Printf("\nindexNumber: %d\n", indexNumber)

	// close the content database creation files ..

	finaliseCollectionDatabase(articleCollectionName, articleJsFile)
	finaliseCollectionDatabase(articleDownloadCollectionName, articleDownloadJsFile)
	finaliseCollectionDatabase(bulletinCollectionName, bulletinJsFile)
	finaliseCollectionDatabase(compendiumDataCollectionName, compendiumDataJsFile)
	finaliseCollectionDatabase(compendiumLandingPageCollectionName, compendiumLandingPageJsFile)
	finaliseCollectionDatabase(datasetLandingPageCollectionName, datasetLandingPageJsFile)
	finaliseCollectionDatabase(staticMethodologyCollectionName, staticMethodologyJsFile)
	finaliseCollectionDatabase(staticMethodologyDownloadCollectionName, staticMethodologyDownloadJsFile)
	finaliseCollectionDatabase(staticQmiCollectionName, staticQmiJsFile)
	finaliseCollectionDatabase(timeseriesCollectionName, timeseriesJsFile)
	finaliseCollectionDatabase(chartCollectionName, chartJsFile)
	finaliseCollectionDatabase(productPageCollectionName, productPageJsFile)
	finaliseCollectionDatabase(tableCollectionName, tableJsFile)
	finaliseCollectionDatabase(equationCollectionName, equationJsFile)
	finaliseCollectionDatabase(imageCollectionName, imageJsFile)
	finaliseCollectionDatabase(releaseCollectionName, releaseJsFile)
	finaliseCollectionDatabase(listCollectionName, listJsFile)
	finaliseCollectionDatabase(staticPageCollectionName, staticPageJsFile)
	finaliseCollectionDatabase(staticAdhocCollectionName, staticAdhocJsFile)
	finaliseCollectionDatabase(referenceTablesCollectionName, referenceTablesJsFile)
	finaliseCollectionDatabase(compendiumChapterCollectionName, compendiumChapterJsFile)
	finaliseCollectionDatabase(staticLandingPageCollectionName, staticLandingPageJsFile)
	finaliseCollectionDatabase(staticArticleCollectionName, staticArticleJsFile)
	finaliseCollectionDatabase(datasetCollectionName, datasetJsFile)
	finaliseCollectionDatabase(timeseriesDatasetCollectionName, timeseriesDatasetJsFile)
	finaliseCollectionDatabase(taxonomyLandingPageCollectionName, taxonomyLandingPageJsFile)

	createContentCountsFile()

	// ===

	if indexNumber > 0 {
		var topics = make([]TopicResponseStore, indexNumber+1)
		var content = make([]ContentResponse, indexNumber+1)

		populateTopicAndContentStructs(topics, content)

		createTopicJsScript(topics)

		createContentJsScript(content)
	}

	createBrokenLinkFile()

	createDuplicatesFile()

	createURICollectionNamesFile()

	fmt.Printf("\nmaxDepth: %d\n", maxDepth)

	fmt.Printf("\nattemptedGetCount is: %v\n", attemptedGetCount)

	fmt.Printf("\nLength of contentDuplicateCheck (URI's saved) is: %d\n", len(contentDuplicateCheck))

	fmt.Printf("\nNumber of URI's seen more than once: %d", len(moreThanOnceURI))

	fmt.Printf("\nNumber of URI's not saved with Version Number on end: %d\n", len(skippedVersionURI))

	fmt.Printf("\nNumber of URI's not saved with # (hash) in them: %d\n", len(hashURI))

	fmt.Printf("\nNumber of URI's not saved with ? (question mark) in them: %d\n", len(questionURI))

	fmt.Printf("\nAll Done.\n")
}
