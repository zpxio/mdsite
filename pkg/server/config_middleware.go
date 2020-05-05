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
	"github.com/apex/log"
	"github.com/gin-gonic/gin"
	"github.com/zpxio/mdsite/pkg/config"
)

const contextConfig = "mdsite-config"

func AddContextConfiguration(conf *config.Values) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(contextConfig, conf)
	}
}

func ContextConfig(c *gin.Context) *config.Values {
	vi := c.Value(contextConfig)

	v, ok := vi.(*config.Values)
	if ok {
		return v
	}

	log.Fatalf("No configuration attached to context: %s", c.FullPath())
	return nil
}

func SiteBaseDirectory(c *gin.Context) string {
	conf := ContextConfig(c)

	return conf.SitePath
}
