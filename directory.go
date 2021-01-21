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
)

// Directory represents a Nuxeo directory
type directory struct {
	EntityType    string                 `json:"entity-type"`
	DirectoryName string                 `json:"directoryName"`
	ID            string                 `json:"id"`
	Properties    map[string]interface{} `json:"properties"`
}

// DirectorySet represents a Nuxeo directory set
type directorySet struct {
	Entries []directory `json:"entries"`
}

func (nuxeoClient *nuxeoClient) GetDirectory(directory string) (directorySet, error) {
	uri := nuxeoClient.url + "/api/v1/directory/" + directory

	resp, err := nuxeoClient.client.R().EnableTrace().Get(uri)

	var directorySet directorySet
	err = HandleResponse(err, resp, &directorySet)

	return directorySet, err
}

func (nuxeoClient *nuxeoClient) CreateDirectory(name string, dir directory) (directory, error) {

	uri := nuxeoClient.url + "/api/v1/directory/" + name

	body, err := json.Marshal(dir)

	resp, err := nuxeoClient.client.R().EnableTrace().SetBody(string(body[:])).Post(uri)

	var newDir directory
	err = HandleResponse(err, resp, &newDir)

	return newDir, err
}
