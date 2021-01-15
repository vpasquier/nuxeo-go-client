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
	EntityType  string                 `json:"entity-type"`
	UID         string                 `json:"uid"`
	Path        string                 `json:"path"`
	Type        string                 `json:"type"`
	Name        string                 `json:"name"`
	Properties  map[string]interface{} `json:"properties"`
	nuxeoClient nuxeoClient
}

// Blob
// type Blob struct {
// 	filename string
// 	size int
// 	file
//   }

type recordSet struct {
	Documents        []document `json:"entries"`
	TotalSize        int        `json:"totalSize"`
	CurrentPageIndex int        `json:"currentPageIndex"`
	NumberOfPages    int        `json:"numberOfPages"`
}

func (doc document) FetchChildren() recordSet {
	url := doc.nuxeoClient.url + "/api/v1/path" + doc.Path + "/@children"

	resp, err := doc.nuxeoClient.client.R().EnableTrace().Get(url)
	var recordSet recordSet
	HandleResponse(err, resp, &recordSet)

	for key, entry := range recordSet.Documents {
		_ = key
		entry.nuxeoClient = doc.nuxeoClient
	}

	// TODO reconnect all documents with nuxeoClient

	return recordSet
}
