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
	"github.com/zpxio/mdsite/pkg/site"
	"net/http"
)

func AttachToc(d *Dispatcher) {
	d.engine.GET("/toc", TableOfContents)
}

func TableOfContents(c *gin.Context) {

	c.Header("Content-Type", gin.MIMEHTML)

	err := config.Global().SiteConfig.Global.TocTemplate.Execute(c.Writer, site.Index())
	if err != nil {
		c.Status(http.StatusInternalServerError)
		c.Error(err)
	}

	c.Status(http.StatusOK)
	c.Done()
}
