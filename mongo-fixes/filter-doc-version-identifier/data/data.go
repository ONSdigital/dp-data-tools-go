package data

import "time"

// ListFilters contains a list of filters
type ListFilters struct {
	Items []Filter
}

// Dataset contains the uniique identifiers that make a dataset unique
type Dataset struct {
	ID      string `bson:"id"        json:"id"`
	Edition string `bson:"edition"   json:"edition"`
	Version int    `bson:"version"   json:"version"`
}

// Filter represents a structure for a filter job
type Filter struct {
	Dataset     *Dataset    `bson:"dataset"              json:"dataset"`
	InstanceID  string      `bson:"instance_id"          json:"instance_id"`
	Dimensions  []Dimension `bson:"dimensions,omitempty" json:"dimensions,omitempty"`
	Downloads   *Downloads  `bson:"downloads,omitempty"  json:"downloads,omitempty"`
	Events      Events      `bson:"events,omitempty"     json:"events,omitempty"`
	FilterID    string      `bson:"filter_id"            json:"filter_id,omitempty"`
	State       string      `bson:"state,omitempty"      json:"state,omitempty"`
	Published   bool        `bson:"published,omitempty"  json:"published,omitempty"`
	Links       LinkMap     `bson:"links"                json:"links,omitempty"`
	LastUpdated time.Time   `bson:"last_updated"         json:"last_updated"`
}

// LinkMap contains a named LinkObject for each link to other resources
type LinkMap struct {
	Dimensions      LinkObject `bson:"dimensions"                 json:"dimensions,omitempty"`
	FilterOutput    LinkObject `json:"filter_output,omitempty" json:"filter_output,omitempty"`
	FilterBlueprint LinkObject `bson:"filter_blueprint,omitempty" json:"filter_blueprint,omitempty"`
	Self            LinkObject `bson:"self"                       json:"self,omitempty"`
	Version         LinkObject `bson:"version"                    json:"version,omitempty"`
}

// LinkObject represents a generic structure for all links
type LinkObject struct {
	ID   string `bson:"id,omitempty" json:"id,omitempty"`
	HRef string `bson:"href"         json:"href,omitempty"`
}

// Dimension represents an object containing a list of dimension values and the dimension name
type Dimension struct {
	URL     string   `bson:"dimension_url,omitempty" json:"dimension_url,omitempty"`
	Name    string   `bson:"name"                    json:"name"`
	Options []string `bson:"options"                 json:"options"`
}

// Downloads represents a list of file types possible to download
type Downloads struct {
	CSV  DownloadItem `bson:"csv"  json:"csv"`
	JSON DownloadItem `bson:"json" json:"json"`
	XLS  DownloadItem `bson:"xls"  json:"xls"`
}

// DownloadItem represents an object containing information for the download item
type DownloadItem struct {
	Size string `bson:"size,omitempty" json:"size"`
	URL  string `bson:"url,omitempty"  json:"url"`
}

// Events represents a list of array objects containing event information against the filter job
type Events struct {
	Error []EventItem `bson:"error,omitempty" json:"error,omitempty"`
	Info  []EventItem `bson:"info,omitempty"  json:"info,omitempty"`
}

// EventItem represents an event object containing event information
type EventItem struct {
	Message string `bson:"message" json:"message,omitempty"`
	Time    string `bson:"time"    json:"time,omitempty"`
	Type    string `bson:"type"    json:"type,omitempty"`
}

// DimensionOption represents dimension option information
type DimensionOption struct {
	DimensionOptionURL string `json:"dimension_option_url"`
	Option             string `json:"option"`
}
