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
	"github.com/zpxio/mdsite/pkg/resource"
	"os"
	"path"
)

var resourceRenderer = make(map[string]resource.Renderer)
var missingRenderer = resource.MissingResource{}

func init() {
	resourceRenderer["md"] = &resource.MarkdownResource{}
}

type PageData struct {
	content []byte
}

func AttachPageHandler(d *Dispatcher) {
	d.engine.NoRoute(Page)
}

func Page(c *gin.Context) {
	resource := c.Request.URL.Path
	renderer, rcFile := FindResourceFile(c, resource)

	renderer.Render(c, rcFile)
}

func FindResourceFile(c *gin.Context, resource string) (resource.Renderer, string) {
	base := SiteBaseDirectory(c)
	rcPrefix := path.Join(base, resource)

	for suffix, renderer := range resourceRenderer {
		rcPath := rcPrefix + "." + suffix
		if fileExists(rcPath) {
			return renderer, rcPath
		}
	}

	return &missingRenderer, resource
}

func fileExists(path string) bool {
	_, err := os.Stat(path)

	if err != nil {
		if os.IsNotExist(err) {
			return false
		}

		log.Errorf("Could not resolve file due to unexpected error: %s", err)
		return false
	} else {
		return true
	}
}

func RenderPath(path string, c *gin.Context) {

}
