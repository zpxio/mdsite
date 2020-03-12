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
	"context"
	"fmt"
	"github.com/apex/log"
	"github.com/gin-gonic/gin"
	"github.com/zpxio/mdsite/pkg/config"
	"net"
	"net/http"
	"net/url"
)

type Dispatcher struct {
	engine     *gin.Engine
	conf       *config.Values
	bindAddr   *net.TCPAddr
	clientAddr string
	server     *http.Server
}

func CreateDispatcher(v *config.Values) *Dispatcher {

	// Create the Gin Engine
	gin.SetMode(gin.ReleaseMode)
	e := gin.Default()
	log.Infof("Gin startup complete")

	d := Dispatcher{
		engine: e,
		conf:   v,
	}

	d.AttachUtility()

	return &d
}

func (d *Dispatcher) AttachUtility() {
	AttachPing(d)
}

func (d *Dispatcher) AttachMiddleware() {

}

func (d *Dispatcher) Start() {

	listenAddr := fmt.Sprintf("%s:%d", d.conf.ListenIp.String(), d.conf.ListenPort)
	log.Infof("Starting server on: %s", listenAddr)

	// Start a listener
	l, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("Failed to set up server socket: %s", err)
	}
	lAddr := l.Addr()
	d.clientAddr = lAddr.String()

	var addrValid bool
	d.bindAddr, addrValid = lAddr.(*net.TCPAddr)
	if !addrValid {
		log.Fatalf("Abnormal binding issue: Listener address is not a TCP address (%T)", lAddr)
	}

	// Report address binding
	listenIp := d.bindAddr.IP.String()
	if listenIp == "::" {
		listenIp = "<all>"
	}
	log.Infof("Listening on interface: %s", listenIp)
	log.Infof("Listening on port: %d", d.bindAddr.Port)

	go func() {
		server := http.Server{Handler: d.engine}
		d.server = &server
		err := d.server.Serve(l)
		if err != nil {
			log.Fatalf("Error while trying to start server: %s", err)
		}
	}()
}

func (d *Dispatcher) Shutdown() {
	err := d.server.Shutdown(context.Background())
	if err != nil {
		log.Warnf("Error while trying to shut down: %s", err)
	}
}

func (d *Dispatcher) ServerUrl() url.URL {
	return url.URL{
		Scheme: "http",
		Host:   d.clientAddr,
	}
}
