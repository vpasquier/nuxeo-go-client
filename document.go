package nuxeoclient

import (
	"encoding/json"
	"errors"
	"log"
)

// Document represents a Nuxeo document
type document struct {
	UID         string `json:"uid"`
	Path        string `json:"path"`
	nuxeoClient nuxeoClient
}

type documents struct {
	Entries []document `json:"entries"`
}

func (doc document) FetchChildren() (documents, error) {
	url := doc.nuxeoClient.url + "/api/v1/path" + doc.Path + "/@children"

	resp, err := doc.nuxeoClient.client.R().EnableTrace().Get(url)

	if err != nil {
		log.Printf("%v", err)
		return documents{}, errors.New("Error while fetching document")
	}

	var documents documents
	jsonErr := json.Unmarshal(resp.Body(), &documents)

	if jsonErr != nil {
		log.Panicf("Can't create Nuxeo Client cause %v", jsonErr)
		return documents, errors.New("Unmarshalling issue with the current user response")
	}

	return documents, nil
}