package nuxeoclient

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSmokeClient(t *testing.T) {
	assert := assert.New(t)
	nuxeoClient := NuxeoClient().URL("https://demo.nuxeo.com/nuxeo").Username("Administrator").Password("Administrator").Debug(false).Build()
	currentUser, err := nuxeoClient.Create()
	if err != nil {
		assert.FailNow("Client should be created")
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
	currentUser, err := nuxeoClient.Create()

	if err != nil {
		assert.FailNow("Client should be crea ted")
	}
	assert.Equal("Administrator", currentUser.Username)
	assert.Equal(true, currentUser.IsAdministrator)
}

func TestRepository(t *testing.T) {
	assert := assert.New(t)

	nuxeoClient := NuxeoClient().URL("https://demo.nuxeo.com/nuxeo").Username("Administrator").Password("Administrator").Debug(false).Build()

	nuxeoClient.Create()

	rootDocument, err := nuxeoClient.FetchDocumentRoot()

	if err != nil {
		assert.FailNow("Document should be fetched")
	}

	assert.Equal(rootDocument.Path, "/")

	domain, err := nuxeoClient.FetchDocumentByPath("/default-domain")

	if err != nil {
		assert.FailNow("Document should be fetched")
	}

	assert.Equal(domain.Path, "/default-domain")

	documents, err := domain.FetchChildren()

	if err != nil {
		assert.FailNow("Document should be fetched")
	}

	assert.Equal(len(documents.Entries), 3)
}
