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
		log.Printf("Can't create Nuxeo Client cause %v", jsonErr)
		return documents, errors.New("Unmarshalling issue with the current user response")
	}

	return documents, nil
}
