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

// User structure
type user struct {
	Username        string                 `json:"id"`
	EntityType      string                 `json:"entity-type"`
	IsAdministrator bool                   `json:"isAdministrator"`
	Properties      map[string]interface{} `json:"properties"`
	IsAnonymous     bool                   `json:"isAnonymous"`
}

type userLogged struct {
	Username        string   `json:"username"`
	EntityType      string   `json:"entity-type"`
	IsAdministrator bool     `json:"isAdministrator"`
	Groups          []string `json:"groups"`
}

// Group structure
type group struct {
	Name string `json:"name"`
}

func (nuxeoClient *nuxeoClient) GetUser(username string) (user, error) {
	uri := nuxeoClient.url + "/api/v1/user/" + username

	resp, err := nuxeoClient.client.R().EnableTrace().Get(uri)

	var returnedUser user
	err = HandleResponse(err, resp, &returnedUser)

	return returnedUser, err
}

func (nuxeoClient *nuxeoClient) CreateUser(newUser user) (user, error) {
	uri := nuxeoClient.url + "/api/v1/user"

	body, err := json.Marshal(newUser)

	resp, err := nuxeoClient.client.R().EnableTrace().SetBody(string(body[:])).Post(uri)

	var returnedUser user
	err = HandleResponse(err, resp, &returnedUser)

	return returnedUser, err
}

func (nuxeoClient *nuxeoClient) DeleteUser(username string) error {
	uri := nuxeoClient.url + "/api/v1/user/" + username

	resp, err := nuxeoClient.client.R().EnableTrace().Delete(uri)

	_ = resp

	return err
}
