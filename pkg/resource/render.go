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

package resource

import (
	"github.com/gin-gonic/gin"
	"html/template"
)

type RenderData struct {
	Resource string

	Title string

	Stylesheets []Stylesheet
	Scripts     []Javascript

	StatusCode int
	MediaType  MediaType

	Content template.HTML
}

type Stylesheet struct {
	Url string
}

type Javascript struct {
	Url    string
	weight uint8
}

func InitRenderData(c *gin.Context, resource string) *RenderData {
	pd := RenderData{
		Resource:  resource,
		MediaType: gin.MIMEHTML,
	}

	return &pd
}
