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
	"flag"
	"github.com/apex/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"net"
	"os"
)

const (
	DefaultSitePath        = "."
	DefaultPort     uint16 = 80
)

var DefaultIp net.IP = net.IPv4(0, 0, 0, 0)

type Values struct {
	SitePath   string
	ListenIp   net.IP
	ListenPort uint16
}

func Create() *Values {
	v := Values{}

	v.setupFlags()

	return &v
}

func (v *Values) setupFlags() {
	viper.SetConfigName("config")

	log.Infof("Initializing configuration")

	// Initialize flags
	pflag.StringVar(&v.SitePath, "site", DefaultSitePath, "The path to the directory containing the site to serve")
	pflag.Uint16Var(&v.ListenPort, "port", DefaultPort, "The port for unencrypted connections")
	pflag.IPVar(&v.ListenIp, "listen", DefaultIp, "The host IP to listen on for connections")
}

func (v *Values) Load() {
	v.LoadAll(os.Args[1:])
}

func (v *Values) LoadAll(a []string) {
	// Parse flags
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.CommandLine.Parse(a)
}
