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
	"net"
	"os"
)

const (
	DefaultSitePath          = "."
	DefaultSiteConfig        = "../config"
	DefaultPort       uint16 = 80
)

var DefaultIp net.IP = net.IPv4(0, 0, 0, 0)

type Values struct {
	SitePath   string
	ConfigPath string
	ListenIp   net.IP
	ListenPort uint16

	TestMode bool
}

func Create() *Values {
	v := Values{}

	v.setupFlags()

	return &v
}

func (v *Values) setupFlags() {
	log.Infof("Initializing configuration")

	// Initialize flags
	pflag.StringVar(&v.SitePath, "site", DefaultSitePath, "The path to the directory containing the site to serve")
	pflag.StringVar(&v.ConfigPath, "config", DefaultSiteConfig, "The path to the directory containing the site configuration")
	pflag.Uint16Var(&v.ListenPort, "port", DefaultPort, "The port for unencrypted connections")
	pflag.IPVar(&v.ListenIp, "listen", DefaultIp, "The host IP to listen on for connections")

	pflag.BoolVar(&v.TestMode, "test", false, "Enable testing mode (integration, not unit)")
}

func (v *Values) Load() {
	v.LoadAll(os.Args[1:])
}

func (v *Values) LoadAll(a []string) {
	log.Infof("Loading configuration.")

	// Parse flags
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	err := pflag.CommandLine.Parse(a)
	if err != nil {
		log.Fatalf("Error reading command options: %s", err)
	}

	// Post-processing, overrides, and inference

	// Test Mode enables ephemeral port and so forth
	if v.TestMode {
		log.Info("Enabling test mode.")
		v.ListenPort = 0
	}
}
