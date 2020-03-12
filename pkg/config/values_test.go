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
	"fmt"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/suite"
	"net"
	"os"
	"testing"
)

type ValuesTestSuite struct {
	suite.Suite
}

func (t *ValuesTestSuite) SetupTest() {
	pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ContinueOnError)
}

func loadVarArgs(v *Values, a ...string) {
	v.LoadAll(a)
}

func (t *ValuesTestSuite) TestValueParse_Defaults() {
	v := Create()
	SetupFlags(v)

	loadVarArgs(v)

	t.Equal(DefaultSitePath, v.SitePath)
	t.Equal(DefaultSiteConfig, v.ConfigPath)
	t.Equal(DefaultPort, v.ListenPort)
	t.Equal(DefaultIp, v.ListenIp)
}

func (t *ValuesTestSuite) TestValueParse_Port() {
	v := Create()
	SetupFlags(v)

	var p uint16 = 23456
	loadVarArgs(v, "--port", fmt.Sprintf("%d", p))

	t.Equal(p, v.ListenPort)
}

func (t *ValuesTestSuite) TestValueParse_Interface() {
	v := Create()
	SetupFlags(v)

	var i = net.IPv4(1, 2, 3, 4)
	loadVarArgs(v, "--listen", i.String())

	t.Equal(i, v.ListenIp)
}

func (t *ValuesTestSuite) TestValueParse_Config() {
	v := Create()
	SetupFlags(v)

	var c = "/etc/config"
	loadVarArgs(v, "--config", c)

	t.Equal(c, v.ConfigPath)
}

func TestValueTestSuite(t *testing.T) {
	suite.Run(t, new(ValuesTestSuite))
}
