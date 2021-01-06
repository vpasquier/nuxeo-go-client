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
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSmokeClient(t *testing.T) {
	assert := assert.New(t)
	nuxeoClient := NuxeoClient().URL("https://demo.nuxeo.com/nuxeo").Username("Administrator").Password("Administrator").Debug(false).Build()
	currentUser, err := nuxeoClient.Login()
	if err != nil {
		assert.Fail("call error")
	}
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

	nuxeoClient := NuxeoClient().URL("https://demo.nuxeo.com/nuxeo").Timeout(1).Headers(headers).Cookies(cookies).Username("Administrator").Password("Administrator").Debug(false).Build()
	currentUser, err := nuxeoClient.Login()

	if err != nil {
		assert.Fail("call error")
	}

	assert.Equal("Administrator", currentUser.Username)
	assert.Equal(true, currentUser.IsAdministrator)
}

func TestRepositoryFetch(t *testing.T) {
	assert := assert.New(t)

	nuxeoClient := NuxeoClient().URL("https://demo.nuxeo.com/nuxeo").Username("Administrator").Password("Administrator").Debug(false).Build()

	nuxeoClient.Login()

	rootDocument, err := nuxeoClient.FetchDocumentRoot()

	if err != nil {
		assert.Fail("call error")
	}

	assert.Equal(rootDocument.Path, "/")

	domain, err := nuxeoClient.FetchDocumentByPath("/default-domain")

	if err != nil {
		assert.Fail("call error")
	}

	assert.Equal(domain.Path, "/default-domain")

	documents := domain.FetchChildren()

	assert.Equal(len(documents.Entries), 3)
}

func TestRepositoryCRUD(t *testing.T) {
	assert := assert.New(t)

	nuxeoClient := NuxeoClient().URL("http://localhost:8080/nuxeo").Username("Administrator").Password("Administrator").Debug(false).Schemas([]string{"*"}).Build()

	nuxeoClient.Login()

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

	if err != nil {
		assert.Fail("call error")
	}

	assert.NotNil(newDocument.UID)
	assert.Equal(newDocument.Path, "/nuxeo/new_file_with_go")
	assert.Equal(newDocument.Properties["dc:title"], "New Document")

	newDocument.Properties["dc:title"] = "Document Updated"
	updatedDocument, err := nuxeoClient.UpdateDocument(newDocument)

	if err != nil {
		assert.Fail("call error")
	}

	assert.Equal(updatedDocument.Properties["dc:title"], "Document Updated")

	err = nuxeoClient.DeleteDocument(updatedDocument)
	if err != nil {
		assert.Fail("call error")
	}

	updatedDocument, err = nuxeoClient.FetchDocumentByPath(updatedDocument.Path)
	if err == nil {
		assert.Fail("This document should not be found")
	}
}
