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
	"github.com/gin-gonic/gin"
	"github.com/zpxio/mdsite/pkg/config"
	"net"
)

type Dispatcher struct {
	engine   *gin.Engine
	bindIp   net.IP
	bindPort uint16
}

func CreateDispatcher(v *config.Values) *Dispatcher {

	// Create the Gin Engine
	e := gin.Default()

	d := Dispatcher{
		engine: e,
		bindIp: v.ListenIp,
	}
	return &d
}

func (d *Dispatcher) AttachMiddleware() {

}

func (d *Dispatcher) Start() {
	d.engine.Run()
}
