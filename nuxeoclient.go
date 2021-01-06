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
	"os"

	log "github.com/sirupsen/logrus"
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
	CreateDocument(parentPath string, input document) (document, error)
	UpdateDocument(input document) (document, error)
	DeleteDocument(input document) error
	// Attack(uri string, body []byte, method string) ([]byte, error)
	// AttachBlob(uid string) error
	// BatchUpload() error
	// Automation(op operation) (output, error)
}

func init() {
	lvl, ok := os.LookupEnv("NUXEO_LOG_LEVEL")
	log.SetOutput(os.Stdout)
	// LOG_LEVEL default to info
	if !ok {
		lvl = "info"
	}
	ll, err := log.ParseLevel(lvl)
	if err != nil {
		ll = log.DebugLevel
	}
	log.SetLevel(ll)
}

// Create the client after applying configuration
func (nuxeoClient *nuxeoClient) Login() (user, error) {

	log.Info("Logging in...")

	url := nuxeoClient.url + "/api/v1/automation/login"

	resp, err := nuxeoClient.client.R().EnableTrace().Post(url)

	var currentUser user
	err = HandleResponse(err, resp, &currentUser)

	log.Info("Logged in")

	return currentUser, err
}

func (nuxeoClient *nuxeoClient) FetchDocumentRoot() (document, error) {

	url := nuxeoClient.url + "/api/v1/path//"

	resp, err := nuxeoClient.client.R().EnableTrace().Get(url)

	var currentDoc document
	err = HandleResponse(err, resp, &currentDoc)

	// Attach client to document
	currentDoc.nuxeoClient = *nuxeoClient

	return currentDoc, err
}

func (nuxeoClient *nuxeoClient) FetchDocumentByPath(path string) (document, error) {

	url := nuxeoClient.url + "/api/v1/path" + path

	resp, err := nuxeoClient.client.R().EnableTrace().Get(url)

	var currentDoc document
	err = HandleResponse(err, resp, &currentDoc)

	currentDoc.nuxeoClient = *nuxeoClient

	return currentDoc, err
}

func (nuxeoClient *nuxeoClient) CreateDocument(parentPath string, input document) (document, error) {
	url := nuxeoClient.url + "/api/v1/path" + parentPath

	body, err := json.Marshal(input)

	resp, err := nuxeoClient.client.R().EnableTrace().SetBody(string(body[:])).Post(url)

	var currentDoc document
	err = HandleResponse(err, resp, &currentDoc)

	currentDoc.nuxeoClient = *nuxeoClient

	return currentDoc, err
}

func (nuxeoClient *nuxeoClient) UpdateDocument(input document) (document, error) {
	url := nuxeoClient.url + "/api/v1/path" + input.Path

	body, err := json.Marshal(input)

	resp, err := nuxeoClient.client.R().EnableTrace().SetBody(string(body[:])).Put(url)

	var currentDoc document
	err = HandleResponse(err, resp, &currentDoc)

	currentDoc.nuxeoClient = *nuxeoClient

	return currentDoc, err
}

func (nuxeoClient *nuxeoClient) DeleteDocument(input document) error {
	url := nuxeoClient.url + "/api/v1/path" + input.Path

	resp, err := nuxeoClient.client.R().EnableTrace().Delete(url)

	err = HandleResponse(err, resp, nil)
	return err
}
