// (C) Copyright 2021 Nuxeo (http:nuxeo.com) and others.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// 	http:www.apache.orglicensesLICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Contributors:
// 	Vladimir Pasquier <vpasquier@nuxeo.com>

package nuxeoclient

import (
	"encoding/json"
	"errors"
	"log"
)

const (
	// DefaultURL is the client url if none has been set
	DefaultURL = "http://localhost:8080/nuxeo"
)

// Client interface
type Client interface {
	Login() (user, error)
	FetchDocumentRoot() (document, error)
	FetchDocumentByPath(path string) (document, error)
	CreateDocument(input document) error
	UpdateDocument(input document) error
	DeleteDocument(uid string) error
	// Attack(uri string, body []byte, method string) ([]byte, error)
	// AttachBlob(uid string) error
	// BatchUpload() error
	// Automation(op operation) (output, error)
}

// Create the client after applying configuration
func (nuxeoClient *nuxeoClient) Login() (user, error) {

	url := nuxeoClient.url + "/api/v1/automation/login"

	resp, err := nuxeoClient.client.R().EnableTrace().Post(url)

	if err != nil {
		log.Printf("%v", err)
		return user{}, errors.New("Client cannot be created")
	}

	data := resp.Body()

	if !json.Valid(data) {
		log.Printf("The response is not json validated")
		return user{}, errors.New("Unmarshalling issue with the current user response")
	}

	var currentUser user
	jsonErr := json.Unmarshal(resp.Body(), &currentUser)

	if jsonErr != nil {
		log.Printf("Can't create Nuxeo Client cause %v", jsonErr)
		return user{}, errors.New("Unmarshalling issue with the current user response")
	}

	log.Println("Nuxeo Client Initialized")

	return currentUser, nil
}

func (nuxeoClient *nuxeoClient) FetchDocumentRoot() (document, error) {

	url := nuxeoClient.url + "/api/v1/path//"

	resp, err := nuxeoClient.client.R().EnableTrace().Get(url)

	if err != nil {
		log.Printf("%v", err)
		return document{}, errors.New("Error while fetching document")
	}

	data := resp.Body()

	if !json.Valid(data) {
		log.Printf("The response is not json validated")
		return document{}, errors.New("Unmarshalling issue with the current document response")
	}

	var currentDoc document
	jsonErr := json.Unmarshal(resp.Body(), &currentDoc)

	if jsonErr != nil {
		log.Printf("Can't create Nuxeo Client cause %v", jsonErr)
		return document{}, errors.New("Unmarshalling issue with the current user response")
	}

	currentDoc.nuxeoClient = *nuxeoClient

	return currentDoc, nil
}

func (nuxeoClient *nuxeoClient) FetchDocumentByPath(path string) (document, error) {

	url := nuxeoClient.url + "/api/v1/path" + path

	resp, err := nuxeoClient.client.R().EnableTrace().Get(url)

	if err != nil {
		log.Printf("%v", err)
		return document{}, errors.New("Error while fetching document")
	}

	data := resp.Body()

	if !json.Valid(data) {
		log.Printf("The response is not json validated")
		return document{}, errors.New("Unmarshalling issue with the current document response")
	}

	var currentDoc document
	jsonErr := json.Unmarshal(resp.Body(), &currentDoc)

	if jsonErr != nil {
		log.Printf("Can't create Nuxeo Client cause %v", jsonErr)
		return document{}, errors.New("Unmarshalling issue with the current user response")
	}

	currentDoc.nuxeoClient = *nuxeoClient

	return currentDoc, nil
}

func (nuxeoClient *nuxeoClient) CreateDocument(input document) error {
	return nil
}

func (nuxeoClient *nuxeoClient) UpdateDocument(input document) error {
	return nil
}

func (nuxeoClient *nuxeoClient) DeleteDocument(uid string) error {
	return nil
}
