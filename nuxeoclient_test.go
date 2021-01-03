package nuxeoclient

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSmokeClient(t *testing.T) {
	assert := assert.New(t)
	nuxeoClient := NuxeoClient().URL("http://localhost:8080/nuxeo").Username("Administrator").Password("Administrator").Debug(true).Build()
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

	nuxeoClient := NuxeoClient().URL("http://localhost:8080/nuxeo").Timeout(10).Headers(headers).Cookies(cookies).Username("Administrator").Password("Administrator").Debug(true).Build()
	currentUser, err := nuxeoClient.Create()
	if err != nil {
		assert.FailNow("Client should be created")
	}
	assert.Equal("Administrator", currentUser.Username)

}
