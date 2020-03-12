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
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

type PingResponse struct {
	Timestamp int64
	ClientIp  string
}

func AttachPing(d *Dispatcher) {
	d.engine.GET("/ping", Ping)
}

func Ping(c *gin.Context) {

	r := PingResponse{
		Timestamp: time.Now().UnixNano() / 1000,
		ClientIp:  c.ClientIP(),
	}

	id := c.GetHeader("X-Ping-Id")
	if id == "" {
		id = fmt.Sprintf("%d", r.Timestamp)
	}

	c.Header("X-Ping-Id", id)
	c.YAML(200, r)
}
