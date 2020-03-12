/*
 * Copyright 2020 zpxio (Jeff Sharpe)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package server

import (
	"github.com/gavv/httpexpect"
	"github.com/stretchr/testify/suite"
	"github.com/zpxio/mdsite/pkg/config"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type PingTestSuite struct {
	suite.Suite
	serverUrl  url.URL
	dispatcher *Dispatcher
	testServer *httptest.Server
}

func (t *PingTestSuite) SetupSuite() {
	conf := config.Create()
	conf.EnableTestMode()

	t.dispatcher = CreateDispatcher(conf)

	t.testServer = httptest.NewServer(t.dispatcher.engine)
}

func (t *PingTestSuite) TearDownSuite() {
	t.testServer.Close()
}

func (t *PingTestSuite) TestPing_Basic() {

	e := httpexpect.New(t.T(), t.testServer.URL)

	id := "8675309"

	e.GET("/ping").
		WithHeader("X-Ping-Id", id).
		Expect().Status(http.StatusOK).Header("X-Ping-Id").Equal(id)
}

func TestPingTestSuite(t *testing.T) {
	suite.Run(t, new(PingTestSuite))
}
