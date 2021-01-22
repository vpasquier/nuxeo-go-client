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
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const DEBUG = false

func TestSmokeClient(t *testing.T) {
	assert := assert.New(t)
	nuxeoClient := NuxeoClient().URL("http://localhost:8080/nuxeo").Username("Administrator").Password("Administrator").Debug(DEBUG).Build()
	currentUser, err := nuxeoClient.Login()
	assert.Nil(err)
	assert.Equal("Administrator", currentUser.Username)
}

func TestClientOptions(t *testing.T) {
	assert := assert.New(t)

	headers := map[string]string{}
	headers["content-type"] = "application/json"
	cookies := []*http.Cookie{
		{
			Name:  "default-cookie-1",
			Value: "This is default-cookie 1 value",
		}, {
			Name:  "default-cookie-2",
			Value: "This is default-cookie 2 value",
		},
	}

	nuxeoClient := NuxeoClient().URL("http://localhost:8080/nuxeo").Timeout(1).Headers(headers).Cookies(cookies).Username("Administrator").Password("Administrator").Debug(DEBUG).Build()
	currentUser, err := nuxeoClient.Login()

	assert.Nil(err)

	assert.Equal("Administrator", currentUser.Username)
	assert.True(currentUser.IsAdministrator)
}

func initTest(t *testing.T) (*assert.Assertions, Client) {
	assert := assert.New(t)

	nuxeoClient := NuxeoClient().URL("http://localhost:8080/nuxeo").Username("Administrator").Password("Administrator").Debug(DEBUG).Schemas([]string{"*"}).Build()

	nuxeoClient.Login()
	return assert, nuxeoClient
}

func TestRepositoryFetch(t *testing.T) {
	assert, nuxeoClient := initTest(t)

	rootDocument, err := nuxeoClient.FetchDocumentRoot()

	assert.Nil(err)

	assert.Equal("/", rootDocument.Path)

	domain, err := nuxeoClient.FetchDocumentByPath("/default-domain")

	assert.Nil(err)

	assert.Equal("/default-domain", domain.Path)

	documents := domain.FetchChildren()

	assert.Equal(3, len(documents.Documents))
}

func TestAsyncFunctions(t *testing.T) {
	assert, nuxeoClient := initTest(t)

	c := make(chan document, 1)

	go nuxeoClient.AsyncFetchDocumentByPath("/default-domain", c)

	select {
	case rootDocument := <-c:
		assert.Equal("/default-domain", rootDocument.Path)
	case <-time.After(1 * time.Second):
		assert.Fail("Result should have been received already")
	}
}

func TestRepositoryCRUD(t *testing.T) {
	assert, nuxeoClient := initTest(t)

	workspaces, err := nuxeoClient.FetchDocumentByPath("/nuxeo")

	if err != nil {
		assert.Fail("call error")
	}

	properties := map[string]interface{}{
		"dc:title": "New Document",
	}

	newDocument := document{
		EntityType: "document",
		Type:       "File",
		Name:       "new_file_with_go",
		Properties: properties,
	}

	newDocument, err = nuxeoClient.CreateDocument(workspaces.Path, newDocument)

	assert.Nil(err)

	assert.NotNil(newDocument.UID)
	assert.Equal("/nuxeo/new_file_with_go", newDocument.Path)
	assert.Equal("New Document", newDocument.Properties["dc:title"])

	newDocument.Properties["dc:title"] = "Document Updated"
	updatedDocument, err := nuxeoClient.UpdateDocument(newDocument)

	assert.Nil(err)

	assert.Equal("Document Updated", updatedDocument.Properties["dc:title"])

	err = nuxeoClient.DeleteDocument(updatedDocument)
	assert.Nil(err)

	updatedDocument, err = nuxeoClient.FetchDocumentByPath(updatedDocument.Path)
	if err == nil {
		assert.Fail("This document should not be found")
	}
}

func TestRepositoryQuery(t *testing.T) {
	assert, nuxeoClient := initTest(t)

	resultSet, err := nuxeoClient.Query("SELECT * FROM Domain")

	assert.Nil(err)

	assert.Equal(1, len(resultSet.Documents))
}
func TestRepositoryDirectory(t *testing.T) {
	assert, nuxeoClient := initTest(t)

	directorySet, err := nuxeoClient.GetDirectory("continent")

	assert.Nil(err)

	assert.Equal(7, len(directorySet.Entries))

	properties := make(map[string]interface{})
	properties["id"] = "go"
	properties["obsolete"] = "0"
	properties["ordering"] = "10"
	properties["label"] = "Go"

	newDir := directory{
		EntityType:    "directoryEntry",
		DirectoryName: "continent",
		Properties:    properties,
	}

	returnedDir, err := nuxeoClient.CreateDirectory("continent", newDir)

	assert.Nil(err)
	assert.NotEmpty(returnedDir.ID)
}

func TestAutomation(t *testing.T) {
	assert, nuxeoClient := initTest(t)

	params := make(map[string]string)

	params["value"] = "/"

	doc, err := nuxeoClient.Automation().Operation("Repository.GetDocument").Parameters(params).DocExecute()

	assert.Nil(err)
	assert.Equal("/", doc.Path)

	params["query"] = "SELECT * FROM Document"
	records, err := nuxeoClient.Automation().Operation("Repository.Query").Parameters(params).DocListExecute()

	assert.Nil(err)
	assert.NotEmpty(records.Documents)

	params["document"] = "/default-domain/workspaces/workspace/file"
	params["save"] = "true"
	params["xpath"] = "file:content"

	image, _ := ioutil.ReadFile("pink.jpg")

	blob, blobError := nuxeoClient.Automation().Operation("Blob.AttachOnDocument").Parameters(params).Blob("pink.jpg", image).BlobExecute()

	assert.Nil(blobError)
	assert.Equal(1025580, len(blob))
}

func TestFetchBlob(t *testing.T) {
	assert, nuxeoClient := initTest(t)

	file, err := nuxeoClient.FetchDocumentByPath("/default-domain/workspaces/workspace/file")

	assert.Nil(err)

	blob, blobError := file.FetchBlob("file:content")

	assert.Nil(blobError)
	assert.Equal(1025580, len(blob))
}

func TestAsyncFetchBlob(t *testing.T) {
	assert, nuxeoClient := initTest(t)

	file, err := nuxeoClient.FetchDocumentByPath("/default-domain/workspaces/workspace/file")

	assert.Nil(err)

	c := make(chan []byte, 1)

	go file.AsyncFetchBlob("file:content", c)

	select {
	case blob := <-c:
		assert.Equal(1025580, len(blob))
	case <-time.After(10 * time.Second):
		assert.Fail("Result should have been received already")
	}
}
