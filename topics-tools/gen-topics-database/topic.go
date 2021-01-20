package main

// TopicResponseStore represents an evolving topic with the current topic and the updated topic
// The 'Next' is what gets updated throughout the publishing journey, and then the 'publish' step copies
// the 'Next' over the 'Current' document, so that 'Current' is whats always returned in the web view.
// NOTE: This TopicResponseStore is slightly different to TopicResponse struct in dp-topic-api for the creation process.
type TopicResponseStore struct {
	ID      string      `bson:"id,omitempty"       json:"id,omitempty"`
	Next    *TopicStore `bson:"next,omitempty"     json:"next,omitempty"`
	Current *TopicStore `bson:"current,omitempty"  json:"current,omitempty"`
}

// TopicStore represents topic schema as it is stored in mongoDB
// and is used for marshaling and unmarshaling json representation for API
// ID is a duplicate of ID in TopicResponseStore, to facilitate each subdocument being a full-formed
// response in its own right depending upon request being in publish or web and also authentication.
// Subtopics contains TopicResonse ID(s).
// NOTE: This TopicStore is slightly different to Topic struct in dp-topic-api for the creation process.
type TopicStore struct {
	ID          string            `bson:"_id,omitempty"            json:"id,omitempty"`
	Description string            `bson:"description,omitempty"    json:"description,omitempty"`
	Title       string            `bson:"title,omitempty"          json:"title,omitempty"`
	Keywords    *[]string         `bson:"keywords,omitempty"       json:"keywords,omitempty"`
	State       string            `bson:"state,omitempty"          json:"state,omitempty"`
	Links       *TopicLinks       `bson:"links,omitempty"          json:"links,omitempty"`
	SubtopicIds *[]string         `bson:"subtopics_ids,omitempty"  json:"subtopics_ids,omitempty"`
	Spotlight   *[]TypeLinkObject `bson:"spotlight,omitempty"      json:"spotlight,omitempty"`
}

// LinkObject represents a generic structure for all links
type LinkObject struct {
	HRef string `bson:"href,omitempty"  json:"href,omitempty"`
	ID   string `bson:"id,omitempty"    json:"id,omitempty"`
}

// TopicLinks represents a list of specific links related to the topic resource
type TopicLinks struct {
	Self      *LinkObject `bson:"self,omitempty"       json:"self,omitempty"`
	Subtopics *LinkObject `bson:"subtopics,omitempty"  json:"subtopics,omitempty"`
	Content   *LinkObject `bson:"content,omitempty"    json:"content,omitempty"`
}
