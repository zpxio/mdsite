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

package config

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type GlobalSuite struct {
	suite.Suite
}

func TestGlobalSuite(t *testing.T) {
	suite.Run(t, new(GlobalSuite))
}

func (s *GlobalSuite) TestUnregistered() {
	globalConf = nil

	s.Panics(func() {
		Global()
	})
}

func (s *GlobalSuite) TestNilRegistration() {
	s.Panics(func() {
		SetGlobal(nil)
	})
}

func (s *GlobalSuite) TestRegistration() {
	v := Create()

	var testPort uint16 = 4242
	v.ListenPort = testPort

	SetGlobal(v)

	s.NotNil(Global())
	s.Equal(testPort, Global().ListenPort)
}
