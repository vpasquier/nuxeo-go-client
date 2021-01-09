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

	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

// HandleResponse handle all responses
func HandleResponse(err error, resp *resty.Response, q interface{}) error {

	if err != nil {
		return errors.Unwrap(err)
	}

	// Explore response object
	log.Debug("Response Info:")
	log.Debug("  Error      :", err)
	log.Debug("  Status Code:", resp.StatusCode())
	log.Debug("  Status     :", resp.Status())
	log.Debug("  Proto      :", resp.Proto())
	log.Debug("  Time       :", resp.Time())
	log.Debug("  Received At:", resp.ReceivedAt())
	log.Debug("  Body       :\n", resp)
	log.Debug()

	// Explore trace info
	log.Debug("Request Trace Info:")
	ti := resp.Request.TraceInfo()
	log.Debug("  DNSLookup     :", ti.DNSLookup)
	log.Debug("  ConnTime      :", ti.ConnTime)
	log.Debug("  TCPConnTime   :", ti.TCPConnTime)
	log.Debug("  TLSHandshake  :", ti.TLSHandshake)
	log.Debug("  ServerTime    :", ti.ServerTime)
	log.Debug("  ResponseTime  :", ti.ResponseTime)
	log.Debug("  TotalTime     :", ti.TotalTime)
	log.Debug("  IsConnReused  :", ti.IsConnReused)
	log.Debug("  IsConnWasIdle :", ti.IsConnWasIdle)
	log.Debug("  ConnIdleTime  :", ti.ConnIdleTime)

	data := resp.Body()

	if resp.StatusCode() != 204 && !json.Valid(data) {
		return errors.New("Json response is not valid")
	}

	if resp.StatusCode() == 404 {
		return errors.New("Cannot find resources")
	}

	if q == nil {
		return nil
	}

	jsonErr := json.Unmarshal(data, q)

	if jsonErr != nil {
		return errors.Unwrap(jsonErr)
	}

	return nil
}
