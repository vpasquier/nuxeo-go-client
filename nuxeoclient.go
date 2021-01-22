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
	"net/url"
	"os"
	"strconv"

	"github.com/go-resty/resty/v2"

	log "github.com/sirupsen/logrus"
)

const (
	// DefaultURL is the client url if none has been set
	DefaultURL = "http://localhost:8080/nuxeo"
)

// Client interface
type Client interface {
	Login() (userLogged, error)
	FetchDocumentRoot() (document, error)
	FetchDocumentByPath(path string) (document, error)
	AsyncFetchDocumentByPath(path string, c chan document)
	CreateDocument(parentPath string, input document) (document, error)
	AsyncCreateDocument(parentPath string, input document, c chan document)
	UpdateDocument(input document) (document, error)
	AsyncUpdateDocument(input document, c chan document)
	DeleteDocument(input document) error
	QueryWithParams(query string, pageSize int, currentPageIndex int, offset int, maxResults int, sortBy string, sortOrder string, queryParams string) (recordSet, error)
	Query(query string) (recordSet, error)
	AsyncQuery(query string, c chan recordSet)
	GetDirectory(directory string) (directorySet, error)
	CreateDirectory(directoryName string, dir directory) (directory, error)
	DeleteDirectory(directoryName string, entry string) error
	Attack(uri string, body []byte, method string) ([]byte, error)
	Automation() Automation
	GetUser(username string) (user, error)
	DeleteUser(username string) error
	CreateUser(newUser user) (user, error)
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
func (nuxeoClient *nuxeoClient) Login() (userLogged, error) {

	url := nuxeoClient.url + "/api/v1/automation/login"

	resp, err := nuxeoClient.client.R().EnableTrace().Post(url)

	var currentUser userLogged
	err = HandleResponse(err, resp, &currentUser)

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

func (nuxeoClient *nuxeoClient) AsyncFetchDocumentByPath(path string, c chan document) {

	url := nuxeoClient.url + "/api/v1/path" + path

	resp, err := nuxeoClient.client.R().EnableTrace().Get(url)

	var currentDoc document
	err = HandleResponse(err, resp, &currentDoc)

	if err != nil {
		panic(err)
	}

	currentDoc.nuxeoClient = *nuxeoClient

	c <- currentDoc
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

func (nuxeoClient *nuxeoClient) AsyncCreateDocument(parentPath string, input document, c chan document) {
	url := nuxeoClient.url + "/api/v1/path" + parentPath

	body, err := json.Marshal(input)

	resp, err := nuxeoClient.client.R().EnableTrace().SetBody(string(body[:])).Post(url)

	var currentDoc document
	err = HandleResponse(err, resp, &currentDoc)

	if err != nil {
		panic(err)
	}

	currentDoc.nuxeoClient = *nuxeoClient

	c <- currentDoc
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

func (nuxeoClient *nuxeoClient) AsyncUpdateDocument(input document, c chan document) {
	url := nuxeoClient.url + "/api/v1/path" + input.Path

	body, err := json.Marshal(input)

	resp, err := nuxeoClient.client.R().EnableTrace().SetBody(string(body[:])).Put(url)

	var currentDoc document
	err = HandleResponse(err, resp, &currentDoc)

	if err != nil {
		panic(err)
	}

	currentDoc.nuxeoClient = *nuxeoClient

	c <- currentDoc
}

func (nuxeoClient *nuxeoClient) DeleteDocument(input document) error {
	url := nuxeoClient.url + "/api/v1/path" + input.Path

	resp, err := nuxeoClient.client.R().EnableTrace().Delete(url)

	err = HandleResponse(err, resp, nil)
	return err
}

func (nuxeoClient *nuxeoClient) Query(query string) (recordSet, error) {
	return nuxeoClient.QueryWithParams(query, 0, 0, 0, 0, "", "", "")
}

func (nuxeoClient *nuxeoClient) AsyncQuery(query string, c chan recordSet) {
	recordSet, err := nuxeoClient.QueryWithParams(query, 0, 0, 0, 0, "", "", "")

	if err != nil {
		panic(err)
	}
	c <- recordSet
}

func (nuxeoClient *nuxeoClient) QueryWithParams(query string, pageSize int, currentPageIndex int, offset int, maxResults int, sortBy string, sortOrder string, queryParams string) (recordSet, error) {

	baseURL, err := url.Parse(nuxeoClient.url)

	_ = err

	baseURL.Path += "/api/v1/search/lang/NXQL/execute"

	// Prepare Query Parameters
	params := url.Values{}
	params.Add("query", query)
	params.Add("pageSize", strconv.Itoa(pageSize))
	params.Add("currentPageIndex", strconv.Itoa(currentPageIndex))
	params.Add("offset", strconv.Itoa(offset))
	params.Add("maxResults", strconv.Itoa(maxResults))
	params.Add("sortBy", sortBy)
	if queryParams != "" {
		params.Add("queryParams", queryParams)
	}

	baseURL.RawQuery = params.Encode()

	resp, err := nuxeoClient.client.R().EnableTrace().Get(baseURL.String())

	var recordSet recordSet
	err = HandleResponse(err, resp, &recordSet)

	for key, doc := range recordSet.Documents {
		_ = key
		doc.nuxeoClient = *nuxeoClient
	}

	return recordSet, err
}

func (nuxeoClient *nuxeoClient) Attack(uri string, body []byte, method string) ([]byte, error) {
	var resp *resty.Response
	var err error
	switch method {
	case "get":
		resp, err = nuxeoClient.client.R().EnableTrace().Get(uri)
	case "post":
		resp, err = nuxeoClient.client.R().EnableTrace().SetBody(string(body[:])).Post(uri)
	case "put":
		resp, err = nuxeoClient.client.R().EnableTrace().SetBody(string(body[:])).Put(uri)
	case "delete":
		resp, err = nuxeoClient.client.R().EnableTrace().Delete(uri)
	default:
		resp, err = nuxeoClient.client.R().EnableTrace().Get(uri)
	}
	return resp.Body(), err
}

func (nuxeoClient *nuxeoClient) Automation() Automation {
	return &automation{
		nuxeoClient: nuxeoClient,
	}
}
