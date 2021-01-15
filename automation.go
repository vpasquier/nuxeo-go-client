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
	"net/url"

	"github.com/go-resty/resty/v2"
)

type automation struct {
	operationName string
	parameters    map[string]string
	context       map[string]string
	input         string
	nuxeoClient   *nuxeoClient
}

type opBody struct {
	Context map[string]string `json:"context"`
	Params  map[string]string `json:"params"`
	Input   string            `json:"input"`
}

// Automation is the automation rest api representation
type Automation interface {
	Operation(name string) Automation
	Parameters(parameters map[string]string) Automation
	Input(input string) Automation
	Context(context map[string]string) Automation
	Execute() (*resty.Response, error)
	DocExecute() (document, error)
	DocListExecute() (recordSet, error)
	BlobExecute() ([]document, error)
}

// Operation name setter
func (auto *automation) Operation(name string) Automation {
	auto.operationName = name
	return auto
}

// Context setter
func (auto *automation) Context(context map[string]string) Automation {
	auto.context = context
	return auto
}

// Parameters setter
func (auto *automation) Parameters(parameters map[string]string) Automation {
	auto.parameters = parameters
	return auto
}

// Input setter
func (auto *automation) Input(input string) Automation {
	auto.input = input
	return auto
}

// Execute returns one of the Automation output type
func (auto *automation) Execute() (*resty.Response, error) {
	baseURL, err := url.Parse(auto.nuxeoClient.url)

	_ = err

	if auto.operationName == "" {
		return nil, errors.New("You should set an operation name")
	}

	if auto.context == nil {
		auto.context = make(map[string]string)
	}

	baseURL.Path += "/site/automation/" + auto.operationName

	opBody := &opBody{
		Context: auto.context,
		Params:  auto.parameters,
	}

	body, err := json.Marshal(opBody)

	response, err := auto.nuxeoClient.client.R().EnableTrace().SetBody(string(body[:])).Post(baseURL.String())

	return response, err
}

// DocExecute returns doc from operation rest api
func (auto *automation) DocExecute() (document, error) {
	response, err := auto.Execute()

	if err != nil {
		return document{}, err
	}

	var currentDoc document
	err = HandleResponse(err, response, &currentDoc)

	currentDoc.nuxeoClient = *auto.nuxeoClient

	return currentDoc, err
}

// DocListExecute returns doc list from operation rest api
func (auto *automation) DocListExecute() (recordSet, error) {
	response, err := auto.Execute()

	if err != nil {
		return recordSet{}, err
	}

	var records recordSet
	err = HandleResponse(err, response, &records)

	for key, entry := range records.Documents {
		_ = key
		entry.nuxeoClient = *auto.nuxeoClient
	}

	return records, err
}

// BlobExecute returns blob from operation rest api
func (auto *automation) BlobExecute() ([]document, error) {
	response, err := auto.Execute()

	if err != nil {
		return []document{}, err
	}

	var docList []document
	err = HandleResponse(err, response, &docList)

	for key, doc := range docList {
		_ = key
		doc.nuxeoClient = *auto.nuxeoClient
	}

	return docList, err
}
