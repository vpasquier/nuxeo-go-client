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

// Document represents a Nuxeo document
type document struct {
	EntityType  string            `json:"entity-type"`
	UID         string            `json:"uid"`
	Path        string            `json:"path"`
	Type        string            `json:"type"`
	Name        string            `json:"name"`
	Properties  map[string]string `json:"properties"`
	nuxeoClient nuxeoClient
}

type documents struct {
	Entries []document `json:"entries"`
}

func (doc document) FetchChildren() documents {
	url := doc.nuxeoClient.url + "/api/v1/path" + doc.Path + "/@children"

	resp, err := doc.nuxeoClient.client.R().EnableTrace().Get(url)
	var documents documents
	HandleResponse(err, resp, &documents)

	return documents
}
