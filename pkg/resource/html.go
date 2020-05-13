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
	"github.com/apex/log"
	"io"
	"io/ioutil"
)

type HtmlResource struct {
	RawResource
}

func (r HtmlResource) MediaType() string {
	return "text/html"
}

func (r HtmlResource) ResourceMode() string {
	return "html"
}

func (r HtmlResource) Render(w io.Writer, data *RenderData) error {
	htData, err := ioutil.ReadFile(data.Resource)

	if err != nil {
		log.Errorf("Failed to read file data [%s]: %s", data.Resource, err)
		return err
	}

	_, err = w.Write(htData)
	if err != nil {
		log.Errorf("Failed to write html data [%s]: %s", data.Resource, err)
		return err
	}

	return nil
}
